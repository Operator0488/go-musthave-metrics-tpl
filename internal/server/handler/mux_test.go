package handler

import (
	"context"
	"errors"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRoute(t *testing.T) {
	svc := &mock.MockService{}
	h := NewStorageHandler(context.Background(), svc)

	router := NewRoute(h)

	req := httptest.NewRequest(
		http.MethodPost,
		"/update/gauge/test/42",
		nil,
	)
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
			h := NewStorageHandler(context.Background(), srv)

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

func Test_getUpdateRequest(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "ok gauge",
			path:    "/update/gauge/test/123",
			wantErr: false,
		},
		{
			name:    "wrong prefix",
			path:    "/wrong/gauge/test/123",
			wantErr: true,
		},
		{
			name:    "wrong type",
			path:    "/update/unknown/test/123",
			wantErr: true,
		},
		{
			name:    "not enough parts",
			path:    "/update/gauge/test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := postUpdateRequest(tt.path)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && req == nil {
				t.Fatal("request is nil")
			}
		})
	}
}
