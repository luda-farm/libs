package gcp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	"cloud.google.com/go/cloudtasks/apiv2beta3/cloudtaskspb"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeEventListenerConfig struct {
	GcpProject          string
	GcpLocation         string
	GcpServiceAccount   string
	StripeWebhookSecret string
}

// Returns a handler that creates cloudtasks from Stripe events.
// The tasks target "POST /stripe/:event_type/:event_subtype/..."
// and the task body contains the event data.
func StripeEventListener(config StripeEventListenerConfig) (http.HandlerFunc, error) {
	client, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("creating cloudtasks client: %w", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		event, err := webhook.ConstructEvent(
			body,
			r.Header.Get("stripe-signature"),
			config.StripeWebhookSecret,
		)
		if err != nil {
			panic(err)
		}

		url := "https://" + r.Host + "/stripe/" + strings.Join(strings.Split(event.Type, "."), "/")
		task := cloudtaskspb.CreateTaskRequest{
			Parent: fmt.Sprintf(
				"projects/%s/locations/%s/queues/stripe-events",
				config.GcpProject, config.GcpLocation,
			),
			Task: &cloudtaskspb.Task{
				PayloadType: &cloudtaskspb.Task_HttpRequest{
					HttpRequest: &cloudtaskspb.HttpRequest{
						HttpMethod: cloudtaskspb.HttpMethod_POST,
						Url:        url,
						Body:       event.Data.Raw,
						AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
							OidcToken: &cloudtaskspb.OidcToken{
								ServiceAccountEmail: config.GcpServiceAccount,
							},
						},
					},
				},
			},
		}

		_, err = client.CreateTask(r.Context(), &task)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}

	return handler, nil
}
