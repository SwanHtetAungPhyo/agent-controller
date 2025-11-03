package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"stock-agent.io/configs"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/handlers/users"
	"stock-agent.io/internal/middleware"
)

type HTTPServer struct {
	router            *gin.Engine
	server            *http.Server
	cfg               *configs.AppConfig
	middlewareManager *middleware.Manager
	store             db.Store
	natsConn          *nats.Conn
}

func NewHTTPServer(
	cfg *configs.AppConfig,
	middlewareManager *middleware.Manager,
	store db.Store,
	natsConn *nats.Conn,
) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	s := &HTTPServer{
		router:            router,
		cfg:               cfg,
		middlewareManager: middlewareManager,
		store:             store,
		natsConn:          natsConn,
	}

	s.setupMiddleware()
	s.setupHealthRoutes()

	return s
}

func (s *HTTPServer) setupMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %d %v \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)
	}))

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     s.cfg.CORSAllowOrigins,
		AllowMethods:     s.cfg.CORSAllowMethods,
		AllowHeaders:     s.cfg.CORSAllowHeaders,
		AllowCredentials: s.cfg.CORSAllowCredentials,
		MaxAge:           time.Duration(s.cfg.CORSMaxAge) * time.Second,
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Total-Count"},
		AllowOriginFunc: func(origin string) bool {
			for _, allowedOrigin := range s.cfg.CORSAllowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					return true
				}
			}
			return false
		},
	}))
}

func (s *HTTPServer) setupHealthRoutes() {
	s.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.router.GET("/ready", func(c *gin.Context) {
		if s.natsConn != nil && s.natsConn.IsConnected() {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		}
	})
}

func RegisterRoutes(
	server *HTTPServer,
	userHandler *users.Handler,
) {
	userHandler.RegisterRoutes(server.router)
}

func (s *HTTPServer) Start(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// HTTP server
			s.server = &http.Server{
				Addr:    ":" + strconv.Itoa(s.cfg.ServerPort),
				Handler: s.router,
			}

			// HTTPS server
			httpsServer := &http.Server{
				Addr:    ":8443",
				Handler: s.router,
			}

			// Start HTTP server
			go func() {
				log.Info().Msgf("Starting HTTP server on port %d", s.cfg.ServerPort)
				if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("Failed to start HTTP server")
				}
			}()

			// Start HTTPS server
			go func() {
				log.Info().Msg("Starting HTTPS server on port 8443")
				if err := httpsServer.ListenAndServeTLS("/certs/core-api.pem", "/certs/core-api-key.pem"); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("Failed to start HTTPS server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down HTTP server...")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := s.server.Shutdown(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Error during server shutdown")
				return err
			}

			if s.natsConn != nil {
				s.natsConn.Drain()
				s.natsConn.Close()
			}

			log.Info().Msg("HTTP server stopped gracefully")
			return nil
		},
	})
}
