package logger

import "go.uber.org/zap"

var instance Logger

type logger struct {
	*zap.Logger
}

type Logger interface {
	Info(msg string, fields ...Field)
}

func InitLogger(lvl string) error {
	level, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = level

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	instance = &logger{zl}

	return nil
}

func Info(msg string, fields ...Field) {
	instance.Info(msg, fields...)
}

func (log *logger) Info(msg string, fields ...Field) {
	log.WithOptions().Info(msg, fields...)
}
