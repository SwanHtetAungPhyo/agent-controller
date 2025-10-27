package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nats-io/nats.go"
	"stock-agent.io/internal/email"
	"stock-agent.io/internal/events"
)

type Server struct {
	NatsConn     *nats.Conn
	EmailService email.EmailService
	EventService *events.EventService
	HealthAddr   string
	udpConn      *net.UDPConn
}

type ServerConfig struct {
	NatsURL       string
	HealthUDPAddr string
	EmailConfig   email.Config
}

func NewServer(cfg ServerConfig) (*Server, error) {
	nc, err := connectNATS(cfg.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	emailService := email.NewEmailService(cfg.EmailConfig)

	eventService := events.NewEventService(nc)

	server := &Server{
		NatsConn:     nc,
		EmailService: *emailService,
		EventService: eventService,
		HealthAddr:   cfg.HealthUDPAddr,
	}

	return server, nil
}

func connectNATS(url string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("stock-agent-server"),
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1), // Infinite reconnects
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Println("NATS connection closed")
		}),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to NATS at %s", nc.ConnectedUrl())
	return nc, nil
}

func (s *Server) StartUDPHealthCheck(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", s.HealthAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}

	s.udpConn = conn
	log.Printf("UDP health check listening on %s", s.HealthAddr)

	go s.handleUDPHealthChecks(ctx)

	return nil
}

func (s *Server) handleUDPHealthChecks(ctx context.Context) {
	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping UDP health check handler")
			return
		default:
			s.udpConn.SetReadDeadline(time.Now().Add(1 * time.Second))

			n, addr, err := s.udpConn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Timeout is expected, continue loop
				}
				log.Printf("UDP read error: %v", err)
				continue
			}

			request := string(buffer[:n])
			response := s.getHealthStatus(request)

			_, err = s.udpConn.WriteToUDP([]byte(response), addr)
			if err != nil {
				log.Printf("UDP write error: %v", err)
			}
		}
	}
}

func (s *Server) getHealthStatus(request string) string {
	status := map[string]string{
		"status": "healthy",
		"nats":   "unknown",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}

	if s.NatsConn != nil && s.NatsConn.IsConnected() {
		status["nats"] = "connected"
	} else {
		status["nats"] = "disconnected"
		status["status"] = "unhealthy"
	}

	switch request {
	case "ping", "health", "status":
		if status["status"] == "healthy" {
			return "OK"
		}
		return "UNHEALTHY"
	case "detailed":
		return fmt.Sprintf("Status: %s, NATS: %s, Time: %s",
			status["status"], status["nats"], status["time"])
	default:
		return "OK"
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.StartUDPHealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to start UDP health check: %w", err)
	}

	if err := s.EventService.Start(ctx); err != nil {
		return fmt.Errorf("failed to start event service: %w", err)
	}

	log.Println("Server started successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	if s.udpConn != nil {
		if err := s.udpConn.Close(); err != nil {
			log.Printf("Error closing UDP connection: %v", err)
		}
	}

	if s.EventService != nil {
		if err := s.EventService.Stop(ctx); err != nil {
			log.Printf("Error stopping event service: %v", err)
		}
	}

	if s.NatsConn != nil {
		s.NatsConn.Drain()
		s.NatsConn.Close()
	}

	log.Println("Server shutdown complete")
	return nil
}

func (s *Server) IsHealthy() bool {
	return s.NatsConn != nil && s.NatsConn.IsConnected()
}
