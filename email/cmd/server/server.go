package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"stock-agent.io/internal/email"
	"stock-agent.io/internal/events"
)

type Server struct {
	NatsConn     *nats.Conn
	EmailService *email.EmailService
	EventService *events.EventService
	HealthAddr   string
	udpConn      *net.UDPConn
	httpServer   *http.Server
	httpsServer  *http.Server
	router       *gin.Engine
}

type ServerConfig struct {
	NatsURL           string
	NatsMaxReconnect  int
	NatsReconnectWait string
	NatsTimeout       string
	HealthUDPAddr     string
	HTTPPort          int
	HTTPSPort         int
	EmailConfig       email.Config
}

func NewServer(cfg ServerConfig) (*Server, error) {
	nc, err := connectNATS(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	emailService, err := email.NewEmailService(cfg.EmailConfig, nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	eventService := events.NewEventService(nc)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		NatsConn:     nc,
		EmailService: emailService,
		EventService: eventService,
		HealthAddr:   cfg.HealthUDPAddr,
		router:       router,
	}

	server.setupRoutes()
	server.setupHTTPServers(cfg)

	return server, nil
}

func connectNATS(cfg ServerConfig) (*nats.Conn, error) {
	timeout, err := time.ParseDuration(cfg.NatsTimeout)
	if err != nil {
		timeout = 10 * time.Second
	}

	reconnectWait, err := time.ParseDuration(cfg.NatsReconnectWait)
	if err != nil {
		reconnectWait = 2 * time.Second
	}

	maxReconnects := cfg.NatsMaxReconnect
	if maxReconnects <= 0 {
		maxReconnects = -1
	}

	opts := []nats.Option{
		nats.Name("stock-agent-server"),
		nats.Timeout(timeout),
		nats.ReconnectWait(reconnectWait),
		nats.MaxReconnects(maxReconnects),
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

	nc, err := nats.Connect(cfg.NatsURL, opts...)
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

	if err := s.startHTTPServers(ctx); err != nil {
		return fmt.Errorf("failed to start HTTP servers: %w", err)
	}

	if err := s.EventService.Start(ctx); err != nil {
		return fmt.Errorf("failed to start event service: %w", err)
	}

	if err := s.EmailService.StartListener(ctx, "email.send"); err != nil {
		return fmt.Errorf("failed to start email listener: %w", err)
	}

	log.Println("Server started successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// Shutdown HTTP servers
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down HTTP server: %v", err)
		}
	}

	if s.httpsServer != nil {
		if err := s.httpsServer.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down HTTPS server: %v", err)
		}
	}

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

func (s *Server) setupRoutes() {
	// Health check endpoints
	s.router.GET("/health", func(c *gin.Context) {
		if s.IsHealthy() {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		}
	})

	s.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Email service endpoints
	api := s.router.Group("/api/v1")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service": "email-service",
				"status":  "running",
				"nats":    s.NatsConn.IsConnected(),
			})
		})

		// Test email endpoint
		api.POST("/send-test-email", func(c *gin.Context) {
			var request struct {
				To      string `json:"to" binding:"required,email"`
				Subject string `json:"subject"`
				Name    string `json:"name"`
				Type    string `json:"type"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set defaults
			if request.Subject == "" {
				request.Subject = "Test Email from Kainos"
			}
			if request.Name == "" {
				request.Name = "Test User"
			}
			if request.Type == "" {
				request.Type = "welcome"
			}

			// Create email event
			emailEvent := map[string]interface{}{
				"type":    request.Type,
				"message": "This is a test email sent through the Kainos email service API.",
				"info": map[string]interface{}{
					"to":      request.To,
					"name":    request.Name,
					"subject": request.Subject,
				},
			}

			// Publish to NATS
			if err := s.EventService.PublishAsync("email.send", &events.Event{
				Type:   "email.send",
				Source: "email-api",
				Data:   emailEvent,
			}); err != nil {
				log.Printf("Failed to publish email event: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to send email",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Email sent successfully",
				"to":      request.To,
				"subject": request.Subject,
				"type":    request.Type,
			})
		})
	}
}

func (s *Server) setupHTTPServers(cfg ServerConfig) {
	httpPort := cfg.HTTPPort
	if httpPort == 0 {
		httpPort = 8082
	}

	httpsPort := cfg.HTTPSPort
	if httpsPort == 0 {
		httpsPort = 8444
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: s.router,
	}

	s.httpsServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpsPort),
		Handler: s.router,
	}
}

func (s *Server) startHTTPServers(ctx context.Context) error {
	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on port %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server
	go func() {
		log.Printf("Starting HTTPS server on port %s", s.httpsServer.Addr)
		if err := s.httpsServer.ListenAndServeTLS("/certs/email-service.pem", "/certs/email-service-key.pem"); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTPS server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) IsHealthy() bool {
	return s.NatsConn != nil && s.NatsConn.IsConnected()
}
