package server

import (
	"crypto/tls"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"stock-agent.io/internal/execution/activities"
	"stock-agent.io/internal/execution/worker"
	work "stock-agent.io/internal/execution/workflow"
)

func (s *HttpServer) temporalWorkerSetup() {
	s.worker = worker.New(
		s.cfg,
		s.temporalClient,
		"default",
	)

	s.worker.RegisterWorkflow(s.workflowManager.StockSummeryWorkflow)

}
func (s *HttpServer) activityManagerSetup() {
	s.activityManager = activities.NewManager(
		s.circuitBreakerClient,
		s.store,
		s.cfg,
	)
}
func (s *HttpServer) temporalSetup() {
	var clientOptions client.Options

	clientOptions = client.Options{
		HostPort:  s.cfg.TemporalHostPort,
		Namespace: s.cfg.TemporalNamespace,
	}

	if s.cfg.TemporalTLS {
		clientOptions.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	}

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Temporal client")
	}

	s.temporalClient = temporalClient
	log.Info().
		Str("host", s.cfg.TemporalHostPort).
		Str("namespace", s.cfg.TemporalNamespace).
		Bool("tls", s.cfg.TemporalTLS).
		Msg("Temporal client connected")
	s.workflowManager = work.NewManager()

}
