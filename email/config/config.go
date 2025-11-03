package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	NatsURL           string `env:"NATS_URL,required"`
	NatsMaxReconnect  int    `env:"NATS_MAX_RECONNECT,required"`
	NatsReconnectWait string `env:"NATS_RECONNECT_WAIT,required"`
	NatsTimeout       string `env:"NATS_TIMEOUT,required"`
	Topic             string `env:"TOPIC,required"`
	ResendAPIKey      string `env:"RESEND_API_KEY,required"`
	FromEmail         string `env:"FROM_EMAIL,required"`
	FromName          string `env:"FROM_NAME,required"`
}

func NewConfig() *Config {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
		return nil
	}
	return cfg
}
