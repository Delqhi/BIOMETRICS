package utils

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	level string
	env   string
	mu    sync.RWMutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func NewLogger(level, env string) *Logger {
	once.Do(func() {
		defaultLogger = &Logger{}
	})

	logger := &Logger{
		level: level,
		env:   env,
	}

	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(parseLevel(level))
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level = zap.NewAtomicLevelAt(parseLevel(level))
	}

	log, err := config.Build()
	if err != nil {
		log = zap.NewNop()
	}

	logger.SugaredLogger = log.Sugar()

	return logger
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) With(fields ...interface{}) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newLogger := &Logger{
		level: l.level,
		env:   l.env,
	}
	newLogger.SugaredLogger = l.SugaredLogger.With(fields...)
	return newLogger
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.SugaredLogger.Fatalw(msg, fields...)
	os.Exit(1)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.SugaredLogger.Errorw(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.SugaredLogger.Warnw(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.SugaredLogger.Infow(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.SugaredLogger.Debugw(msg, fields...)
}

func Default() *Logger {
	if defaultLogger == nil {
		return NewLogger("info", "development")
	}
	return defaultLogger
}
