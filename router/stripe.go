package router

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/luda-farm/libs/errorutil"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type (
	StripeEventHandlingConfig struct {
		GcpProject, GcpLocation, StripeWebhookSecret string
	}
)

func (router *Router) InitStripeEventHandling(config StripeEventHandlingConfig) {
	client := errorutil.Must(cloudtasks.NewClient(context.Background()))
	router.Handle(http.MethodPost, "/stripe/events",
		func(ctx Context) {
			requestBody := errorutil.Must(io.ReadAll(ctx.Request.Body))
			event, err := webhook.ConstructEvent(
				requestBody,
				ctx.Request.Header.Get("stripe-signature"),
				config.StripeWebhookSecret,
			)
			ctx.CheckClientError(err, http.StatusUnauthorized)
			url := "https://" + ctx.Request.Host + eventTypeToResource(event.Type)
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
			errorutil.Must(client.CreateTask(context.Background(), &task))
			ctx.Response.WriteHeader(http.StatusNoContent)
		},
	)
}

func eventTypeToResource(event string) string {
	return "/stripe/" + strings.Join(strings.Split(event, "."), "/")
}

func (router *Router) HandleStripeEvent(event string, handler Handler) {
	router.Handle(http.MethodPost, eventTypeToResource(event), handler)
}
