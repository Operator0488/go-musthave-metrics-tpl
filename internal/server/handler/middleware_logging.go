package handler

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriterLog struct {
	http.ResponseWriter
	status int
	size   int
}

func logging(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriterLog{
				ResponseWriter: w,
				status:         0,
				size:           0,
			}

			next.ServeHTTP(rw, r)
			duration := time.Since(start)
			log.Info("HTTP запрос",
				zap.String("URL", r.URL.Path),
				zap.String("method", r.Method),
				zap.Duration("time", duration),
			)

			log.Info("HTTP ответ",
				zap.Int("status", rw.status),
				zap.Int("size", rw.size),
			)
		})
	}
}

func (rw *responseWriterLog) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterLog) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}
