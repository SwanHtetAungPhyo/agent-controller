package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *HttpServer) Stop() {
	serverCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Stopping the temporal Worker....")
	s.worker.Stop()

	if s.temporalClient != nil {
		s.temporalClient.Close()
		log.Info().Msg("Temporal client closed...")
	}

	if s.databasePool != nil {
		s.databasePool.Close()
		log.Info().Msg("Database pool closed...")
	}

	log.Info().Msg("Server shutting down...")
	if err := s.server.Shutdown(serverCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}
}
