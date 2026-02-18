package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"biometrics/internal/api"
	"biometrics/internal/cache"
	"biometrics/internal/config"
	"biometrics/internal/database"
	"biometrics/pkg/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := utils.NewLogger(cfg.LogLevel, cfg.Environment)
	logger.Info("Starting biometrics API server",
		"port", cfg.Server.Port,
		"environment", cfg.Environment,
	)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal("Database connection failed", "error", err)
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		logger.Fatal("Redis connection failed", "error", err)
	}
	defer redisClient.Close()
	logger.Info("Redis connected successfully")

	router := api.SetupRouter(db.DB, redisClient, logger, cfg)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	go func() {
		logger.Info("Server listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited properly")
}
