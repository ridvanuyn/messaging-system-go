package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ridvanuyn/messaging-system-go/internal/api"
	"github.com/ridvanuyn/messaging-system-go/internal/config"
	"github.com/ridvanuyn/messaging-system-go/internal/repository"
	"github.com/ridvanuyn/messaging-system-go/internal/service"
	"github.com/ridvanuyn/messaging-system-go/internal/worker"
	"github.com/ridvanuyn/messaging-system-go/pkg/database"
	_ "github.com/ridvanuyn/messaging-system-go/docs" 

)

// @title Messaging API
// @version 1.0
// @description API for automatic message sending system
// @host localhost:8080
// @BasePath /api
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Database connection
	db, err := database.NewPostgresDB(cfg.DbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Redis connection
	redis, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Repository
	messageRepo := repository.NewMessageRepository(db, redis)

	// Service
	messageService := service.NewMessageService(messageRepo, cfg)

	// Scheduler
	scheduler := worker.NewScheduler(messageService)

	// HTTP Handler
	handler := api.NewHandler(messageService, scheduler)

	// Router
	router := api.SetupRouter(handler)

	// HTTP Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Server running at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Start scheduler automatically on deployment
	scheduler.Start()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Stop scheduler
	scheduler.Stop()

	// Shut down server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
