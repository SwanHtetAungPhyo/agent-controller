package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"stock-agent.io/configs"
)

func NewNATSConnection(cfg *configs.AppConfig) (*nats.Conn, error) {
	if cfg.NATSUrl == "" {
		return nil, fmt.Errorf("NATS URL is required")
	}

	opts := []nats.Option{
		nats.Name("kainos-core-api"),
		nats.Timeout(30 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Error().Err(err).Msg("NATS disconnected")
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info().Str("url", nc.ConnectedUrl()).Msg("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Info().Msg("NATS connection closed")
		}),
	}

	// Retry connection with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		log.Info().
			Str("url", cfg.NATSUrl).
			Int("attempt", i+1).
			Int("max_retries", maxRetries).
			Msg("Attempting to connect to NATS")

		nc, err := nats.Connect(cfg.NATSUrl, opts...)
		if err == nil {
			log.Info().Str("url", nc.ConnectedUrl()).Msg("Connected to NATS")
			return nc, nil
		}

		log.Warn().
			Err(err).
			Int("attempt", i+1).
			Msg("Failed to connect to NATS, retrying...")

		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 2 * time.Second
			log.Info().Dur("wait_time", waitTime).Msg("Waiting before retry")
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("failed to connect to NATS after %d attempts", maxRetries)
}
