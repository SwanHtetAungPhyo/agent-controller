package server

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	db "stock-agent.io/db/sqlc"
)

func (s *HttpServer) DatabaseSetup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	s.databasePool, err = pgxpool.New(ctx, s.cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	err = s.databasePool.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
		return
	}
	log.Info().Msg("Connected to database")
	s.store = db.NewStore(s.databasePool)

}
