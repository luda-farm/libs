package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type (
	StripeEventRouter struct {
		router *Router
	}
	StripeEventRouterConfig struct {
		Host, GcpProject, GcpLocation, StripeWebhookSecret string
	}
)

func NewStripeEventRouter(router *Router, config StripeEventRouterConfig) StripeEventRouter {
	client, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	router.Post("/stripe/events", func(ctx Context) {
		payload := ctx.RawRequestBody()
		signature := ctx.Request.Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, signature, config.StripeWebhookSecret)
		if err != nil {
			ctx.Response.WriteHeader(http.StatusUnauthorized)
			return
		}

		task := tasks.CreateTaskRequest{
			Parent: fmt.Sprintf(
				"projects/%s/locations/%s/queues/stripe-events",
				config.GcpProject, config.GcpLocation,
			),
			Task: &tasks.Task{
				MessageType: &tasks.Task_HttpRequest{
					HttpRequest: &tasks.HttpRequest{
						HttpMethod: tasks.HttpMethod_POST,
						Url:        config.Host + pathFromEvent(event.Type),
						Body:       event.Data.Raw,
					},
				},
			},
		}

		if _, err := client.CreateTask(context.Background(), &task); err != nil {
			panic(err)
		}

		ctx.Response.WriteHeader(http.StatusNoContent)
	})

	return StripeEventRouter{
		router: router,
	}
}

func pathFromEvent(event string) string {
	return "/stripe/" + strings.Join(strings.Split(event, "."), "/")
}

func (router StripeEventRouter) Handle(event string, handler Handler) {
	router.router.Post(pathFromEvent(event), handler)
}
