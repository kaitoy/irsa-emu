// The irsa-emu-webhook is a mutating admission webhook server
// that insert a sidecar into a pod with IRSA's role annotation.
package main

import (
	"flag"

	"github.com/kaitoy/irsa-emu/pkg/logging"
	"github.com/kaitoy/irsa-emu/pkg/webhook"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		bindAddr = flag.String(
			"bind-addr",
			":443",
			"the address where the HTTPS server will be listening to serve the webhooks.",
		)
		tlsCert          = flag.String("tls-cert", "", "the path for the webhook HTTPS server TLS cert file.")
		tlsKey           = flag.String("tls-key", "", "the path for the webhook HTTPS server TLS key file.")
		sidecarImageRepo = flag.String(
			"sidecar-image-repo",
			"kaitoy/irsa-emu-creds-injector",
			"the repository of the sidecar image.",
		)
		sidecarImageTag = flag.String(
			"sidecar-image-tag",
			"latest",
			"the tag of the sidecar image.",
		)
		awsEnvvarsSecret = flag.String(
			"aws-envvars-secret",
			"irsa-emu-aws-envvars",
			"the name of the secret containing AWS envvars for the sidecar.",
		)
		stsEndpointURL = flag.String(
			"sts-endpoint-url",
			"",
			"the endpoint URL of STS that the sidecar accesses.",
		)
	)
	flag.Parse()

	logging.Init(logrus.DebugLevel)
	logger := logging.GetLogger()

	logger.Infof(
		"Starting webhook server. "+
			"bindAddr: %s, tlsCert: %s, tlsKey: %s, "+
			"sidecarImageRepo: %s, sidecarImageTag: %s, awsEnvvarsSecret: %s, stsEndpointURL: %s",
		*bindAddr, *tlsCert, *tlsKey, *sidecarImageRepo, *sidecarImageTag, *awsEnvvarsSecret, *stsEndpointURL,
	)
	if err := webhook.NewServer(
		*bindAddr, *tlsCert, *tlsKey, *sidecarImageRepo, *sidecarImageTag, *awsEnvvarsSecret, *stsEndpointURL,
	).Run(); err != nil {
		logger.Errorf("An error occurred: %v", err)
	}
}
