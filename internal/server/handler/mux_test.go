package handler

import (
	"bytes"
	"errors"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRoute(t *testing.T) {
	svc := &mock.MockService{}
	h := NewStorageHandler(svc)
	router := NewChiRoute(h)

	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	v := []byte(`{
  "id": "PollCount4",
  "type": "gauge",
  "value": 1.0007
	} `)

	req := httptest.NewRequest(
		http.MethodPost,
		"/update",
		bytes.NewReader(v),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !svc.Called {
		t.Fatal("service was not called")
	}
}

func TestStorageHandler_PostUpdate(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		serviceErr error
		wantCode   int
	}{
		{
			name:     "ok",
			path:     "/update/gauge/test/10",
			wantCode: http.StatusOK,
		},
		{
			name:     "bad path",
			path:     "/bad/path",
			wantCode: http.StatusBadRequest,
		},
		{
			name:       "service error",
			path:       "/update/gauge/test/10",
			serviceErr: errors.New("service failed"),
			wantCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &mock.MockService{
				Err: tt.serviceErr,
			}
			h := NewStorageHandler(srv)

			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rec := httptest.NewRecorder()

			h.PostUpdate(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("expected %d, got %d", tt.wantCode, rec.Code)
			}
		})
	}
}

func Test_checkPost(t *testing.T) {
	handler := checkPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("", "")
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name     string
		method   string
		wantCode int
	}{
		{"POST ok", http.MethodPost, http.StatusOK},
		{"GET bad", http.MethodGet, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/any", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("expected %d, got %d", tt.wantCode, rec.Code)
			}
		})
	}
}
