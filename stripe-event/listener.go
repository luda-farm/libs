package stripeevent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/luda-farm/libs/std"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type EventListenerBuilder struct {
	GcpProject          string
	GcpLocation         string
	GcpServiceAccount   string
	OidcAudience        string
	StripeWebhookSecret string
}

// Returns a handler that creates cloudtasks from Stripe events.
// The tasks target "POST /stripe/:event_type/:event_subtype/..."
// and the task body contains the event data.
func (config EventListenerBuilder) Build() http.HandlerFunc {
	client := std.Must(cloudtasks.NewClient(context.Background()))
	return func(w http.ResponseWriter, r *http.Request) {
		event := std.Must(webhook.ConstructEvent(
			std.Must(io.ReadAll(r.Body)),
			r.Header.Get("stripe-signature"),
			config.StripeWebhookSecret,
		))

		url := "https://" + r.Host + "/stripe/" + strings.Join(strings.Split(event.Type, "."), "/")
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
						AuthorizationHeader: &tasks.HttpRequest_OidcToken{
							OidcToken: &tasks.OidcToken{
								Audience:            config.OidcAudience,
								ServiceAccountEmail: config.GcpServiceAccount,
							},
						},
					},
				},
			},
		}
		std.Must(client.CreateTask(context.Background(), &task))
		w.WriteHeader(http.StatusNoContent)
	}
}
