package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"myapp/handlers"
	"myapp/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using defaults.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Router with middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.HealthCheck)

	handler := middleware.Logger(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server running on http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown: %s\n", err)
	}
}
