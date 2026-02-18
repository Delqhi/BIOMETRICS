package workers

import (
	"biometrics/internal/config"
	"biometrics/pkg/utils"

	"context"
)

type Queue struct {
	log *utils.Logger
}

func NewQueue(cfg config.RedisConfig, logger *utils.Logger) (*Queue, error) {
	return &Queue{log: logger}, nil
}

type CaptchaSolverWorker struct {
	log *utils.Logger
}

func NewCaptchaSolverWorker(db interface{}, queue *Queue, logger *utils.Logger, cfg config.CaptchaConfig) *CaptchaSolverWorker {
	return &CaptchaSolverWorker{log: logger}
}

func (w *CaptchaSolverWorker) Start(ctx context.Context) {}

type SurveyWorker struct {
	log *utils.Logger
}

func NewSurveyWorker(db interface{}, queue *Queue, logger *utils.Logger, cfg config.SurveyConfig) *SurveyWorker {
	return &SurveyWorker{log: logger}
}

func (w *SurveyWorker) Start(ctx context.Context) {}
