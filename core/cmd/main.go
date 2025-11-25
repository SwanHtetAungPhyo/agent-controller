package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"stock-agent.io/internal/database"
	fxModules "stock-agent.io/internal/fx"
	"stock-agent.io/internal/server"
	"stock-agent.io/internal/temporal"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Use console output only in containerized environment
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	app := fx.New(
		fxModules.ConfigModule,
		database.DatabaseModule(),
		fxModules.NATSModule,
		fxModules.EventsModule,
		temporal.TemporalModule(),
		fxModules.MiddlewareModule,
		fxModules.HandlersModule,
		fxModules.ServerModule,
		fx.Invoke(func(server *server.HTTPServer, lc fx.Lifecycle) {
			server.Start(lc)
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed to start application")
	}

	log.Info().Msg("Application started successfully")

	<-app.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("Error during application shutdown")
	}

	log.Info().Msg("Application stopped gracefully")
}
