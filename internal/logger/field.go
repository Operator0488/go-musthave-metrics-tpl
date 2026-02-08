package logger

import (
	"go.uber.org/zap"
	"time"
)

type Field = zap.Field

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Error(val error) Field {
	return zap.Error(val)
}
