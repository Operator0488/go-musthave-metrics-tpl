package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapLogger struct {
	l *zap.Logger
}

//go:generate mockgen -source=logger.go -destination=./mocks/mock_logger.go -package=mocks
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

func New(lvl string) (*ZapLogger, error) {
	level, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания логгера: %w", err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = level

	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания логгера: %w", err)
	}

	return &ZapLogger{l: zl}, nil
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field) {
	zl.l.Info(msg, fields...)
}

func (zl *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{
		zl.l.With(fields...),
	}
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field) {
	zl.l.Error(msg, fields...)
}
