package temporal

import (
	"context"
	"crypto/tls"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
	"stock-agent.io/configs"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/execution/activities"
	"stock-agent.io/internal/execution/worker"
	"stock-agent.io/internal/execution/workflow"
	"stock-agent.io/pkg/circuitBreaker"
)

func NewTemporalClient(cfg *configs.AppConfig) (client.Client, error) {
	clientOptions := client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.TemporalNamespace,
	}

	if cfg.TemporalTLS {
		clientOptions.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	}

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("host", cfg.TemporalHostPort).
		Str("namespace", cfg.TemporalNamespace).
		Bool("tls", cfg.TemporalTLS).
		Msg("Temporal client connected")

	return temporalClient, nil
}

func NewWorkflowManager() *workflow.Manager {
	return workflow.NewManager()
}

func NewActivityManager(circuitBreakerClient *circuitBreaker.Client, store db.Store, cfg *configs.AppConfig) *activities.Manager {
	return activities.NewManager(circuitBreakerClient, store, cfg)
}

func NewWorker(cfg *configs.AppConfig, temporalClient client.Client) *worker.Worker {
	return worker.New(cfg, temporalClient, "default")
}

func NewCircuitBreakerClient() *circuitBreaker.Client {
	return &circuitBreaker.Client{}
}

func NewScheduleClient(temporalClient client.Client) client.ScheduleClient {
	return temporalClient.ScheduleClient()
}

func TemporalModule() fx.Option {
	return fx.Module("temporal",
		fx.Provide(
			NewTemporalClient,
			NewWorkflowManager,
			NewActivityManager,
			NewWorker,
			NewCircuitBreakerClient,
			NewScheduleClient,
		),
		fx.Invoke(func(lc fx.Lifecycle, temporalClient client.Client, worker *worker.Worker) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := worker.Start(); err != nil {
							log.Error().Err(err).Msg("Failed to start Temporal worker")
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info().Msg("Stopping Temporal worker")
					worker.Stop()
					temporalClient.Close()
					return nil
				},
			})
		}),
	)
}
