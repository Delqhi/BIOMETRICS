package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"biometrics-cli/internal/collision"
	"biometrics-cli/internal/opencode"
	"biometrics-cli/internal/project"
	"biometrics-cli/internal/prompt"
	"biometrics-cli/internal/quality"
	"biometrics-cli/internal/telemetry"
)

func main() {
	logger := telemetry.SetupLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutdown signal received")
		cancel()
	}()

	logger.Info("=== BIOMETRICS 24/7 ULTRA ORCHESTRATOR STARTING ===")

	modelPool := collision.NewModelPool()
	executor := opencode.NewExecutor(logger)
	basePath := "/Users/jeremy/.sisyphus/plans"

	for {
		if ctx.Err() != nil {
			break
		}

		projects, err := project.DiscoverProjects(basePath)
		if err != nil || len(projects) == 0 {
			logger.Warn("No projects found, waiting...")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, projID := range projects {
			cycleCtx, traceID := telemetry.InjectTraceID(ctx)
			logger.Info("Processing project", slog.String("project", projID), slog.String("trace_id", traceID))

			task, err := project.GetNextTask(projID)
			if err != nil {
				logger.Debug("No pending tasks", slog.String("project", projID))
				continue
			}

			logger.Info("Executing Task", slog.String("task_id", task.ID))

			model := "qwen-3.5"
			modelPool.Acquire(cycleCtx, model)

			taskPrompt := prompt.GenerateEnterprisePrompt(projID, "active_plan.md", task.ID, task.Description)

			req := opencode.AgentRequest{
				ProjectID: projID,
				Model:     model,
				Prompt:    taskPrompt,
				Category:  "build",
			}

			result := executor.RunAgent(cycleCtx, req)
			modelPool.Release(model)

			if !result.Success {
				logger.Error("Task failed", slog.String("task_id", task.ID), slog.String("error", result.Error.Error()))
				continue
			}

			if err := quality.EnforceQualityGate(cycleCtx, executor, req); err == nil {
				project.MarkTaskCompleted(projID, task.ID)
				logger.Info("Task completed and verified", slog.String("task_id", task.ID))
			} else {
				logger.Error("Quality Gate failed", slog.String("task_id", task.ID))
			}
		}

		time.Sleep(5 * time.Second)
	}
}
