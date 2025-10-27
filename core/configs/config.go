package configs

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	// App Configuration
	AppName        string `env:"APP_NAME" envDefault:""`
	AppEnvironment string `env:"APP_APP_ENVIRONMENT" envDefault:"development"`
	AppDebug       bool   `env:"APP_APP_DEBUG" envDefault:"false"`

	// Server Configuration
	ServerHost string `env:"APP_SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort int    `env:"APP_SERVER_PORT" envDefault:"8081"`

	// Authentication
	ClerkSecret string `env:"CLERK_SECRET" envDefault:""`
	JWTSecret   string `env:"APP_JWT_SECRET" envDefault:""`
	JWTTTL      int    `env:"APP_JWT_TTL" envDefault:"8640"`

	// Database Configuration
	DatabaseHost     string `env:"APP_DATABASE_HOST" envDefault:"localhost"`
	DatabasePort     int    `env:"APP_DATABASE_PORT" envDefault:"5432"`
	DatabaseUsername string `env:"APP_DATABASE_USERNAME" envDefault:""`
	DatabasePassword string `env:"APP_DATABASE_PASSWORD" envDefault:""`
	DatabaseName     string `env:"APP_DATABASE_NAME" envDefault:""`
	DatabaseSSLMode  string `env:"APP_DATABASE_SSL_MODE" envDefault:"disable"`
	DatabaseURL      string `env:"DATABASE_URL" envDefault:""`

	// Redis Configuration
	RedisHost     string `env:"APP_REDIS_HOST" envDefault:"localhost"`
	RedisPort     int    `env:"APP_REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"APP_REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"APP_REDIS_DB" envDefault:"0"`

	// Temporal Configuration
	TemporalHostPort  string `env:"APP_TEMPORAL_HOSTPORT" envDefault:"localhost:7233"`
	TemporalNamespace string `env:"APP_TEMPORAL_NAMESPACE" envDefault:"default"`
	TemporalTLS       bool   `env:"APP_TEMPORAL_TLS" envDefault:"false"`

	// Svix Configuration
	SvixSecret string `env:"APP_SVIX_SECRET" envDefault:""`
	SvixAppID  string `env:"APP_SVIX_APP_ID" envDefault:""`

	// Core API
	CoreAPIGRPCEndpoint string `env:"APP_CORE_API_GRPC_ENDPOINT" envDefault:""`

	// Version Information
	PostgreSQLVersion string `env:"POSTGRESQL_VERSION" envDefault:""`
	RedisVersion      string `env:"REDIS_VERSION" envDefault:""`
	TemporalVersion   string `env:"TEMPORAL_VERSION" envDefault:""`
	TemporalUIVersion string `env:"TEMPORAL_UI_VERSION" envDefault:""`
	MastraUrl         string `env:"MASTRA_URL" envDefault:""`

	// CORS Configuration
	CORSAllowOrigins     []string `env:"CORS_ALLOW_ORIGINS" envSeparator:"," envDefault:"*"`
	CORSAllowMethods     []string `env:"CORS_ALLOW_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"`
	CORSAllowHeaders     []string `env:"CORS_ALLOW_HEADERS" envSeparator:"," envDefault:"Origin,Content-Length,Content-Type,Authorization"`
	CORSAllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	CORSMaxAge           int      `env:"CORS_MAX_AGE" envDefault:"43200"` // 12 hours in seconds
}

// NewAppConfig creates and loads application configuration
func NewAppConfig() *AppConfig {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}

// loadConfig loads configuration from environment variables
func loadConfig() (*AppConfig, error) {
	// Load .env file if it exists (optional)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &AppConfig{}

	// Parse environment variables into struct
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate performs validation on the configuration
func (c *AppConfig) Validate() error {
	// Add custom validation logic here
	// Example:
	// if c.JWTSecret == "" {
	//     return fmt.Errorf("JWT_SECRET is required")
	// }
	return nil
}

// IsDevelopment checks if the app is running in development mode
func (c *AppConfig) IsDevelopment() bool {
	return c.AppEnvironment == "development"
}

// IsProduction checks if the app is running in production mode
func (c *AppConfig) IsProduction() bool {
	return c.AppEnvironment == "production"
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *AppConfig) GetDatabaseDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUsername,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}

// GetRedisAddr returns the Redis connection address
func (c *AppConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}
