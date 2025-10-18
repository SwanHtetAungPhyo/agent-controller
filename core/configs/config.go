package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName        string
	AppEnvironment string
	AppDebug       bool

	ServerHost           string
	ServerPort           int
	ServerAllowedOrigins []string // Add this field

	ClerkSecret string

	JWTSecret string
	JWTTTL    int

	DatabaseHost     string
	DatabasePort     int
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
	DatabaseURL      string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	TemporalHostPort  string
	TemporalNamespace string
	TemporalTLS       bool

	SvixSecret string
	SvixAppID  string

	CoreAPIGRPCEndpoint string

	PostgreSQLVersion string
	RedisVersion      string
	TemporalVersion   string
	TemporalUIVersion string
	MastraUrl         string

	// CORS specific fields
	CORSAllowOrigins     []string
	CORSAllowMethods     []string
	CORSAllowHeaders     []string
	CORSAllowCredentials bool
	CORSMaxAge           int
}

func NewAppConfig() *AppConfig {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}

func loadConfig() (*AppConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &AppConfig{
		AppName:             getEnv("APP_NAME", ""),
		AppEnvironment:      getEnv("APP_APP_ENVIRONMENT", "development"),
		AppDebug:            getEnvBool("APP_APP_DEBUG", false),
		ServerHost:          getEnv("APP_SERVER_HOST", "0.0.0.0"),
		ServerPort:          getEnvInt("APP_SERVER_PORT", 8081),
		ClerkSecret:         getEnv("CLERK_SECRET", ""),
		JWTSecret:           getEnv("APP_JWT_SECRET", ""),
		JWTTTL:              getEnvInt("APP_JWT_TTL", 8640),
		DatabaseHost:        getEnv("APP_DATABASE_HOST", "localhost"),
		DatabasePort:        getEnvInt("APP_DATABASE_PORT", 5432),
		DatabaseUsername:    getEnv("APP_DATABASE_USERNAME", ""),
		DatabasePassword:    getEnv("APP_DATABASE_PASSWORD", ""),
		DatabaseName:        getEnv("APP_DATABASE_NAME", ""),
		DatabaseSSLMode:     getEnv("APP_DATABASE_SSL_MODE", "disable"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		RedisHost:           getEnv("APP_REDIS_HOST", "localhost"),
		RedisPort:           getEnvInt("APP_REDIS_PORT", 6379),
		RedisPassword:       getEnv("APP_REDIS_PASSWORD", ""),
		RedisDB:             getEnvInt("APP_REDIS_DB", 0),
		TemporalHostPort:    getEnv("APP_TEMPORAL_HOSTPORT", "localhost:7233"),
		TemporalNamespace:   getEnv("APP_TEMPORAL_NAMESPACE", "default"),
		TemporalTLS:         getEnvBool("APP_TEMPORAL_TLS", false),
		SvixSecret:          getEnv("APP_SVIX_SECRET", ""),
		SvixAppID:           getEnv("APP_SVIX_APP_ID", ""),
		CoreAPIGRPCEndpoint: getEnv("APP_CORE_API_GRPC_ENDPOINT", ""),
		PostgreSQLVersion:   getEnv("POSTGRESQL_VERSION", ""),
		RedisVersion:        getEnv("REDIS_VERSION", ""),
		TemporalVersion:     getEnv("TEMPORAL_VERSION", ""),
		TemporalUIVersion:   getEnv("TEMPORAL_UI_VERSION", ""),
		MastraUrl:           getEnv("MASTRA_URL", ""),

		// CORS Configuration
		CORSAllowOrigins:     getEnvSlice("CORS_ALLOW_ORIGINS", []string{"*"}, ","),
		CORSAllowMethods:     getEnvSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, ","),
		CORSAllowHeaders:     getEnvSlice("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, ","),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
		CORSMaxAge:           getEnvInt("CORS_MAX_AGE", 12*60*60), // 12 hours in seconds
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "TRUE"
	}
	return defaultValue
}

// getEnvSlice parses environment variable as a slice of strings
func getEnvSlice(key string, defaultValue []string, separator string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}
