package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"stock-agent.io/cmd/server"
	"stock-agent.io/configs"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
func main() {
	cfg := configs.NewAppConfig()
	log.Info().Interface("config", cfg).Msg("Starting application")
	app := server.NewHttpServer(cfg)
	go func() {
		app.Start()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info().Msg("Shutting down...")
	app.Stop()
}
