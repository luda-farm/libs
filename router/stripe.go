package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type (
	StripeEventHandlingConfig struct {
		GcpProject, GcpLocation, StripeWebhookSecret string
	}
)

func (router *Router) InitStripeEventHandling(config StripeEventHandlingConfig) error {
	client, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect to cloudtasks: %w", err)
	}

	router.Handle(http.MethodPost, "/stripe/events",
		func(res http.ResponseWriter, req *http.Request, _ map[string]string) {
			requestBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			event, err := webhook.ConstructEvent(
				requestBody,
				req.Header.Get("stripe-signature"),
				config.StripeWebhookSecret,
			)
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			url := "https://" + req.Host + pathFromEvent(event.Type)
			task := tasks.CreateTaskRequest{
				Parent: fmt.Sprintf(
					"projects/%s/locations/%s/queues/stripe-events",
					config.GcpProject, config.GcpLocation,
				),
				Task: &tasks.Task{
					MessageType: &tasks.Task_HttpRequest{
						HttpRequest: &tasks.HttpRequest{
							HttpMethod: tasks.HttpMethod_POST,
							Url:        url,
							Body:       event.Data.Raw,
						},
					},
				},
			}

			_, err = client.CreateTask(context.Background(), &task)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
			}

			res.WriteHeader(http.StatusNoContent)
		},
	)
	return nil
}

func pathFromEvent(event string) string {
	return "/stripe/" + strings.Join(strings.Split(event, "."), "/")
}

func (router *Router) HandleStripeEvent(event string, handler Handler) {
	router.Handle(http.MethodPost, pathFromEvent(event), handler)
}
