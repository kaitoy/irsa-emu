package webhook

import (
	"fmt"
	"net/http"

	"github.com/kaitoy/irsa-emu/pkg/logging"
	kwhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
)

// Server starts a Webhook server.
type Server interface {
	// Run starts a Webhook server.
	Run() error
}

type server struct {
	bindAddr         string
	tlsCert          string
	tlsKey           string
	sidecarImageRepo string
	sidecarImageTag  string
	awsEnvvarsSecret string
	stsEndpointURL   string
}

// NewServer creates a Server instance.
func NewServer(
	bindAddr string,
	tlsCert string,
	tlsKey string,
	sidecarImageRepo string,
	sidecarImageTag string,
	awsEnvvarsSecret string,
	stsEndpointURL string,
) Server {
	return &server{
		bindAddr,
		tlsCert,
		tlsKey,
		sidecarImageRepo,
		sidecarImageTag,
		awsEnvvarsSecret,
		stsEndpointURL,
	}
}

// Run starts the server.
func (r *server) Run() error {
	wh, err := kwmutating.NewWebhook(kwmutating.WebhookConfig{
		ID:      "irsa-emu",
		Mutator: GetMutatorFunc(r.sidecarImageRepo, r.sidecarImageTag, r.awsEnvvarsSecret, r.stsEndpointURL),
		Logger:  logging.GetLogger(),
	})
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}

	whHandler, err := kwhttp.HandlerFor(kwhttp.HandlerConfig{Webhook: wh, Logger: logging.GetLogger()})
	if err != nil {
		return fmt.Errorf("create webhook handler: %w", err)
	}

	mux := http.NewServeMux()
	addHealthzEndpoint(mux)
	mux.Handle("/mutate", whHandler)

	err = http.ListenAndServeTLS(r.bindAddr, r.tlsCert, r.tlsKey, mux)
	if err != nil {
		return fmt.Errorf("serve webhook: %w", err)
	}

	return nil
}

func addHealthzEndpoint(mux *http.ServeMux) {
	mux.HandleFunc(
		"/healthz",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		}),
	)
}
