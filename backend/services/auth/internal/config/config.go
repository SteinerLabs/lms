package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the configuration for the auth service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port int
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

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  int // in minutes
	RefreshTokenTTL int // in days
	Issuer          string
	Audience        string
}

// KafkaConfig represents the Kafka configuration
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Server configuration
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "50051"))
	if err != nil {
		return nil, fmt.Errorf("invalid server port: %v", err)
	}

	// Database configuration
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid database port: %v", err)
	}

	// JWT configuration
	accessTokenTTL, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_TTL", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token TTL: %v", err)
	}

	refreshTokenTTL, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_TTL", "7"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token TTL: %v", err)
	}

	// Kafka configuration
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")

	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "auth"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:  accessTokenTTL,
			RefreshTokenTTL: refreshTokenTTL,
			Issuer:          getEnv("JWT_ISSUER", "lms-auth-service"),
			Audience:        getEnv("JWT_AUDIENCE", "lms-api"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{kafkaBrokers},
			Topic:   getEnv("KAFKA_TOPIC", "user-events"),
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
