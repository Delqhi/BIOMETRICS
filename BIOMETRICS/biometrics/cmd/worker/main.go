package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"biometrics/internal/config"
	"biometrics/internal/database"
	"biometrics/internal/workers"
	"biometrics/pkg/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := utils.NewLogger(cfg.LogLevel, cfg.Environment)
	logger.Info("Starting biometrics worker",
		"environment", cfg.Environment,
		"worker_count", cfg.Worker.Count,
	)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal("Database connection failed", "error", err)
	}
	defer db.Close()

	queue, err := workers.NewQueue(cfg.Redis, logger)
	if err != nil {
		logger.Fatal("Queue initialization failed", "error", err)
	}

	captchaWorker := workers.NewCaptchaSolverWorker(db, queue, logger, cfg.Worker.Captcha)
	surveyWorker := workers.NewSurveyWorker(db, queue, logger, cfg.Worker.Survey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	surveyWorker.Start(ctx)
	logger.Info("Survey worker started")

	captchaWorker.Start(ctx)
	logger.Info("Captcha solver worker started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down workers...")
	cancel()
	logger.Info("Workers stopped")
}
