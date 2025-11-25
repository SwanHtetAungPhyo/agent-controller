package fx

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"go.uber.org/fx"
	"stock-agent.io/configs"
	"stock-agent.io/internal/events"
	"stock-agent.io/internal/handlers/users"
	"stock-agent.io/internal/handlers/workflow"
	"stock-agent.io/internal/middleware"
	natsClient "stock-agent.io/internal/nats"
	"stock-agent.io/internal/server"
)

var ConfigModule = fx.Module("config",
	fx.Provide(configs.NewAppConfig),
)

var NATSModule = fx.Module("nats",
	fx.Provide(natsClient.NewNATSConnection),
)

var EventsModule = fx.Module("events",
	fx.Provide(events.NewPublisher),
)

var HandlersModule = fx.Module("handlers",
	fx.Provide(
		users.NewHandler,
	),
	fx.Provide(workflow.NewHandler),
)

var MiddlewareModule = fx.Module("middleware",
	fx.Provide(
		func(cfg *configs.AppConfig) string {
			return cfg.ClerkSecret
		},
		func(cfg *configs.AppConfig) *clerk.ClientConfig {
			return &clerk.ClientConfig{
				BackendConfig: clerk.BackendConfig{
					Key: &cfg.ClerkSecret,
				},
			}
		},
		middleware.NewManager,
	),
)

var ServerModule = fx.Module("server",
	fx.Provide(server.NewHTTPServer),
	fx.Invoke(server.RegisterRoutes),
)
