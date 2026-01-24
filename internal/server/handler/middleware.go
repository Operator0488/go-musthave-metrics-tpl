package handler

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/logger"
	"net/http"
	"time"
)

type middleware func(h http.Handler) http.Handler

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func middleConveyor(h http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func checkPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorBadRequest(w, "Invalid Method")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		logger.Info("HTTP запрос",
			logger.String("URL", r.URL.Path),
			logger.String("method", r.Method),
			logger.Duration("time", duration),
		)

		logger.Info("HTTP ответ",
			logger.Int("status", rw.status),
			logger.Int("size", rw.size),
		)
	})
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}
