package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/repository"
	"github.com/SteinerLabs/lms/backend/services/auth/internal/service"
)

func main() {
	fmt.Println("Auth service")
	for {
		time.Sleep(10 * time.Second)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database
	log.Println("Initializing database...")
	err = repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Create the auth service
	authService, err := service.NewAuthServiceImpl(cfg)
	if err != nil {
		log.Fatalf("Failed to create auth service: %v", err)
	}

	// Create a simple HTTP server for health checks
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.Server.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Close any other resources
	if err := authService.Close(); err != nil {
		log.Printf("Error closing auth service: %v", err)
	}

	log.Println("Server exited")
}
