package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"stock-agent.io/configs"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/execution/activities"
	"stock-agent.io/internal/execution/worker"
	work "stock-agent.io/internal/execution/workflow"
	"stock-agent.io/internal/handlers/users"
	"stock-agent.io/internal/handlers/workflow"
	"stock-agent.io/internal/middleware"
	"stock-agent.io/internal/routes"
	"stock-agent.io/pkg/circuitBreaker"
)

type HttpServer struct {
	router               *gin.Engine
	server               *http.Server
	route                *routes.Route
	middlewareManager    *middleware.Manager
	clerkCfg             *clerk.ClientConfig
	databasePool         *pgxpool.Pool
	handlers             handlers
	store                db.Store
	cfg                  *configs.AppConfig
	temporalClient       client.Client
	worker               *worker.Worker
	workflowManager      *work.Manager
	activityManager      *activities.Manager
	scheduleClient       client.ScheduleClient
	circuitBreakerClient *circuitBreaker.Client
}

type handlers struct {
	userHandler *users.Handler
	workflow    *workflow.Handler
}

func NewHttpServer(cfg *configs.AppConfig) *HttpServer {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	clerkSecret := cfg.ClerkSecret
	if clerkSecret == "" {
		log.Fatal().Msg("Clerk secret is required")
	}

	clerkConfig := &clerk.ClientConfig{
		BackendConfig: clerk.BackendConfig{
			HTTPClient:           nil,
			URL:                  nil,
			Key:                  &clerkSecret,
			CustomRequestHeaders: nil,
		},
	}

	server := &HttpServer{
		router:   router,
		clerkCfg: clerkConfig,
		cfg:      cfg,
	}

	clerk.SetKey(clerkSecret)
	server.DatabaseSetup()
	server.temporalSetup()
	server.activityManagerSetup()
	server.temporalWorkerSetup()
	server.handlerSetup()
	server.setupMiddleware()

	return server
}

func (s *HttpServer) setupMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
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
				if allowedOrigin == "*" {
					return true
				}
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
	}))
}

func (s *HttpServer) Start() {
	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.cfg.ServerPort),
		Handler: s.router,
	}

	go func() {
		err := s.worker.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start worker")
			return
		}
	}()
	log.Info().Msgf("Starting server on port %d", s.cfg.ServerPort)
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
