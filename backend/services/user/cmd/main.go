package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/SteinerLabs/lms/backend/services/user/internal/config"
	"github.com/SteinerLabs/lms/backend/shared/events"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
)

//go:embed schema.sql
var schema string

func main() {
	log.Println("User service")
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	con, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	db, err := InitDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	js, err := jetstream.New(con)
	if err != nil {
		panic(err)
	}

	consumer, err := events.NewConsumer(js, db, "TEST", "TEST.test", "user-service")
	if err != nil {
		panic(err)
	}

	err = consumer.Start(
		func(ctx context.Context, event events.Event[any]) error {
			fmt.Printf("Received event: %+v\n", event)
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	done := make(chan bool)
	<-done
	log.Println("Done")
}

func InitDB(cfg *config.Config) (*sql.DB, error) {
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
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Execute the schema
	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return db, nil
}
