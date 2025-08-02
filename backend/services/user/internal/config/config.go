package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the configuration for the auth service
type Config struct {
	Database DatabaseConfig
	Nats     NatsConfig
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type NatsConfig struct {
	URL string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Database configuration
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid database port: %v", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "auth"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Nats: NatsConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
