package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Server ServerConfig

	// gRPC configuration
	GRPC GRPCConfig

	// Service authentication
	ServiceKey string

	// Internal gRPC address for inter-service communication
	InternalGRPCAddr string

	// Notification service configuration
	Notification NotificationConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port string
	Host string
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port        string
	Host        string
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

// NotificationConfig holds notification service configuration
type NotificationConfig struct {
	Provider string // "email", "mock", etc.
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")

	cfg.GRPC.Port = getEnv("GRPC_PORT", "9090")
	cfg.GRPC.Host = getEnv("GRPC_HOST", "0.0.0.0")
	cfg.GRPC.TLSEnabled = getEnvAsBool("GRPC_TLS_ENABLED", true)
	cfg.GRPC.TLSCertFile = getEnv("GRPC_TLS_CERT_FILE", "certs/server.crt")
	cfg.GRPC.TLSKeyFile = getEnv("GRPC_TLS_KEY_FILE", "certs/server.key")

	cfg.ServiceKey = getEnv("SERVICE_KEY", "")
	if cfg.ServiceKey == "" {
		return nil, fmt.Errorf("SERVICE_KEY environment variable is required")
	}

	cfg.InternalGRPCAddr = getEnv("INTERNAL_GRPC_ADDR", fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port))

	cfg.Notification.Provider = getEnv("NOTIFICATION_PROVIDER", "email")

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

