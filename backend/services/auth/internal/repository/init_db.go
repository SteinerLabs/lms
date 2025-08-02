package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDB initializes the database with the schema
func InitDB(cfg *config.Config) error {
	// Create the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Get the schema file path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}
	schemaPath := filepath.Join(filepath.Dir(filename), "schema.sql")

	// Read the schema file
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	schema := string(schemaBytes)

	// Execute the schema
	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
