package logger

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

var errFailedToSetLevel = errors.New("failed to set log level")

func New(level, file string) (*Logger, error) {
	var lvl zapcore.Level

	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("%w: %+v", errFailedToSetLevel, err)
	}

	cfg := zap.Config{
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{file},
		ErrorOutputPaths: []string{file},
		Level:            zap.NewAtomicLevelAt(lvl),
		EncoderConfig:    zapcore.EncoderConfig{MessageKey: "M"},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}
