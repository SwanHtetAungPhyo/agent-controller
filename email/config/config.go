package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	NatUrl string `env:"NAT_URL"`
	Topic  string `env:"TOPIC"`
	Email  string `env:"EMAIL"`
	Passwd string `env:"PASSWD"`
}

func NewConfig() *Config {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
		return nil
	}
	return cfg
}
