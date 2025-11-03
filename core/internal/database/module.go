package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"stock-agent.io/configs"
	db "stock-agent.io/db/sqlc"
)

func NewDatabaseConnection(cfg *configs.AppConfig) (*pgxpool.Pool, error) {
	var databaseURL string
	if cfg.DatabaseURL != "" {
		databaseURL = cfg.DatabaseURL
	} else {
		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.DatabaseUsername,
			cfg.DatabasePassword,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName,
			cfg.DatabaseSSLMode,
		)
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.DatabaseHost).
		Int("port", cfg.DatabasePort).
		Str("database", cfg.DatabaseName).
		Msg("Connected to database")

	return pool, nil
}

func NewStore(pool *pgxpool.Pool) db.Store {
	return db.NewStore(pool)
}

func DatabaseModule() fx.Option {
	return fx.Module("database",
		fx.Provide(
			NewDatabaseConnection,
			NewStore,
		),
		fx.Invoke(func(lc fx.Lifecycle, pool *pgxpool.Pool) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Info().Msg("Closing database connection")
					pool.Close()
					return nil
				},
			})
		}),
	)
}
