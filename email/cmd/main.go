package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"stock-agent.io/cmd/server"
	"stock-agent.io/config"
	"stock-agent.io/internal/email"
	"stock-agent.io/internal/events"
)

func main() {
	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configs := config.NewConfig()
	if configs == nil {
		log.Fatal().Msg("config is nil")
	}

	cfg := server.ServerConfig{
		NatsURL:       configs.NatUrl,
		HealthUDPAddr: getEnv("HEALTH_UDP_ADDR", ":8080"),
		EmailConfig: email.Config{
			SMTPHost:     "",
			SMTPPort:     0,
			SMTPUsername: "",
			SMTPPassword: "",
			FromEmail:    "",
			FromName:     "",
			UseTLS:       false,
			UseSSL:       false,
		},
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create server")
	}

	setupEventHandlers(srv.EventService)

	if err := srv.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Server is running... Press Ctrl+C to stop")

	<-sigChan
	log.Info().Msg("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error during shutdown")
	}

	log.Info().Msg("Server stopped gracefully")
}

func setupEventHandlers(es *events.EventService) {
	SingleErrorMust(es.Subscribe("user.created", func(event *events.Event) error {
		log.Info().
			Str("event_type", "user.created").
			Interface("data", event.Data).
			Msg("User created event received")
		return nil
	}))

	SingleErrorMust(es.Subscribe("user.updated", func(event *events.Event) error {
		log.Info().
			Str("event_type", "user.updated").
			Interface("data", event.Data).
			Msg("User updated event received")
		return nil
	}))

	SingleErrorMust(es.QueueSubscribe("order.process", "order-workers", func(event *events.Event) error {
		log.Info().
			Str("event_type", "order.process").
			Str("queue", "order-workers").
			Interface("data", event.Data).
			Msg("Processing order event")
		return nil
	}))

	SingleErrorMust(es.SubscribeWithReply("stock.quote", func(event *events.Event) (*events.Event, error) {
		symbol := event.Data["symbol"].(string)
		log.Info().
			Str("event_type", "stock.quote").
			Str("symbol", symbol).
			Msg("Stock quote requested")

		response := &events.Event{
			ID:        "response-" + event.ID,
			Type:      "stock.quote.response",
			Timestamp: time.Now().UTC(),
			Source:    "stock-service",
			Data: map[string]interface{}{
				"symbol": symbol,
				"price":  150.25,
				"volume": 1000000,
			},
		}

		return response, nil
	}))

	log.Info().Msg("Event handlers registered")
}

func publishExampleEvents(es *events.EventService) {
	userEvent := &events.Event{
		ID:     "evt-123",
		Type:   "user.created",
		Source: "user-service",
		Data: map[string]interface{}{
			"user_id": "12345",
			"email":   "user@example.com",
			"name":    "John Doe",
		},
	}

	if err := es.Publish("user.created", userEvent); err != nil {
		log.Error().Err(err).Msg("Failed to publish event")
	}

	quoteRequest := &events.Event{
		ID:     "req-456",
		Type:   "stock.quote.request",
		Source: "client",
		Data: map[string]interface{}{
			"symbol": "AAPL",
		},
	}

	response, err := es.Request("stock.quote", quoteRequest, 5*time.Second)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get quote")
	} else {
		log.Info().
			Interface("response_data", response.Data).
			Msg("Received stock quote response")
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func SingleErrorMust(err error) {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func Must[T any](val T, err error) T {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Debug().Interface("value", val).Msg("Must operation succeeded")
	return val
}
