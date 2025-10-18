package server

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
	"stock-agent.io/pkg/circuitBreaker"
)

func (s *HttpServer) setupCircuitBreaker() {
	s.circuitBreakerClient = circuitBreaker.NewCircuitBreakerClient(&circuitBreaker.Config{
		Name:        "Mastra-ai-agent",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     65 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Info().Msgf("Circuit Breaker %s: %s -> %s", name, from, to)
		},
	})
}
