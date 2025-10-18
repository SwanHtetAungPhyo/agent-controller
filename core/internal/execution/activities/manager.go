package activities

import (
	"stock-agent.io/configs"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/pkg/circuitBreaker"
)

type Manager struct {
	circuitBreaker *circuitBreaker.Client
	store          db.Store
	cfg            *configs.AppConfig
}

func NewManager(
	circuitBreaker *circuitBreaker.Client,
	store db.Store,
	cfg *configs.AppConfig,
) *Manager {
	return &Manager{
		circuitBreaker: circuitBreaker,
		store:          store,
		cfg:            cfg,
	}
}
