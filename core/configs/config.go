package configs

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName        string `env:"APP_NAME,required"`
	AppEnvironment string `env:"APP_APP_ENVIRONMENT,required"`
	AppDebug       bool   `env:"APP_APP_DEBUG,required"`

	ServerHost string `env:"APP_SERVER_HOST,required"`
	ServerPort int    `env:"APP_SERVER_PORT,required"`

	ClerkSecret string `env:"CLERK_SECRET,required"`
	JWTSecret   string `env:"APP_JWT_SECRET,required"`
	JWTTTL      int    `env:"APP_JWT_TTL,required"`

	DatabaseHost     string `env:"APP_DATABASE_HOST,required"`
	DatabasePort     int    `env:"APP_DATABASE_PORT,required"`
	DatabaseUsername string `env:"APP_DATABASE_USERNAME,required"`
	DatabasePassword string `env:"APP_DATABASE_PASSWORD,required"`
	DatabaseName     string `env:"APP_DATABASE_NAME,required"`
	DatabaseSSLMode  string `env:"APP_DATABASE_SSL_MODE,required"`
	DatabaseURL      string `env:"DATABASE_URL"`

	RedisHost     string `env:"APP_REDIS_HOST,required"`
	RedisPort     int    `env:"APP_REDIS_PORT,required"`
	RedisPassword string `env:"APP_REDIS_PASSWORD,required"`
	RedisDB       int    `env:"APP_REDIS_DB,required"`

	TemporalHostPort  string `env:"APP_TEMPORAL_HOSTPORT,required"`
	TemporalNamespace string `env:"APP_TEMPORAL_NAMESPACE,required"`
	TemporalTLS       bool   `env:"APP_TEMPORAL_TLS,required"`

	NATSUrl string `env:"APP_NATS_URL,required"`

	SvixSecret string `env:"APP_SVIX_SECRET,required"`
	SvixAppID  string `env:"APP_SVIX_APP_ID,required"`

	CORSAllowOrigins     []string `env:"CORS_ALLOW_ORIGINS,required" envSeparator:","`
	CORSAllowMethods     []string `env:"CORS_ALLOW_METHODS,required" envSeparator:","`
	CORSAllowHeaders     []string `env:"CORS_ALLOW_HEADERS,required" envSeparator:","`
	CORSAllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS,required"`
	CORSMaxAge           int      `env:"CORS_MAX_AGE,required"`
}

func NewAppConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse configuration: %v", err)
	}

	return cfg
}

func (c *AppConfig) IsDevelopment() bool {
	return c.AppEnvironment == "development"
}

func (c *AppConfig) IsProduction() bool {
	return c.AppEnvironment == "production"
}

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

func (c *AppConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}
