package stripewebhook

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/luda-farm/libs/errutil"
	"github.com/stripe/stripe-go/v72/webhook"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// Returns a handler that creates cloudtasks routed to "/stripe/:event_type/:event_subtype/..."
func NewListener(gcpProject, gcpLocation, webhookSecret string) http.HandlerFunc {
	client := errutil.Must(cloudtasks.NewClient(context.Background()))
	return func(w http.ResponseWriter, r *http.Request) {
		event := errutil.Must(webhook.ConstructEvent(
			errutil.Must(io.ReadAll(r.Body)),
			r.Header.Get("stripe-signature"),
			webhookSecret,
		))
		url := "https://" + r.Host + "/stripe/" + strings.Join(strings.Split(event.Type, "."), "/")
		task := tasks.CreateTaskRequest{
			Parent: fmt.Sprintf(
				"projects/%s/locations/%s/queues/stripe-events",
				gcpProject, gcpLocation,
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
		errutil.Must(client.CreateTask(context.Background(), &task))
		w.WriteHeader(http.StatusNoContent)
	}
}
