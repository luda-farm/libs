package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/luda-farm/libs/assert"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type (
	StripeEventHandlingConfig struct {
		GcpProject, GcpLocation, StripeWebhookSecret string
	}
)

func (router *Router) InitStripeEventHandling(config StripeEventHandlingConfig) {
	client := assert.Must(cloudtasks.NewClient(context.Background()))

	router.Post("/stripe/events", func(ctx Context) {
		payload := ctx.ReadBytes()
		signature := ctx.Request.Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, signature, config.StripeWebhookSecret)
		if err != nil {
			ctx.Response.WriteHeader(http.StatusUnauthorized)
			return
		}

		url := "https://" + ctx.Request.Host + pathFromEvent(event.Type)
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

		assert.Must(client.CreateTask(context.Background(), &task))
		ctx.Response.WriteHeader(http.StatusNoContent)
	})
}

func pathFromEvent(event string) string {
	return "/stripe/" + strings.Join(strings.Split(event, "."), "/")
}

func (router *Router) HandleStripeEvent(event string, handler Handler) {
	router.Post(pathFromEvent(event), handler)
}
