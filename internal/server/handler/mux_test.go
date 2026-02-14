package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocklog "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger/mocks"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	mockserv "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestStorageHandler_GetValue_ValueOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	val := 123.45
	resp := models.GetValueResponse{Value: &val}

	svc.EXPECT().
		SenderGetValue(models.GetValueRequest{MType: models.Gauge, ID: "cpu"}).
		Return(resp, nil)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/cpu", nil)
	req = withChiParams(req, map[string]string{
		"type": models.Gauge,
		"name": "cpu",
	})

	w := httptest.NewRecorder()
	h.GetValue(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/html", w.Header().Get("Content-Type"))
	require.Equal(t, "123.45", w.Body.String())
}

func TestStorageHandler_GetValue_BadType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	req := httptest.NewRequest(http.MethodGet, "/value/bad/cpu", nil)
	req = withChiParams(req, map[string]string{
		"type": "bad",
		"name": "cpu",
	})
	w := httptest.NewRecorder()

	h.GetValue(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_GetValues_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)
	mapa := map[string]interface{}{
		"LastGc": models.MetricStore{
			ID:    "LastGc",
			MType: models.Gauge,
			Value: 123.2,
		},
	}

	svc.EXPECT().
		SenderGetValues().
		Return(mapa, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.GetValues(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/html", w.Header().Get("Content-Type"))
	require.Equal(t, fmt.Sprint(mapa), w.Body.String())
}

func TestStorageHandler_PostUpdate_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	svc.EXPECT().
		SenderPostUpdate(models.PostUpdateRequest{
			MType:    models.Counter,
			ID:       "hits",
			ValueStr: "10",
		}).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/update/counter/hits/10", nil)
	w := httptest.NewRecorder()

	h.PostUpdate(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestStorageHandler_PostUpdateWithBody_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString("{bad json"))
	w := httptest.NewRecorder()

	h.PostUpdateWithBody(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_PostUpdateWithBody_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	val := 99.9

	body := models.PostUpdateRequest{
		MType: models.Gauge,
		ID:    "ram",
		Value: &(val),
	}
	b, _ := json.Marshal(body)

	svc.EXPECT().
		SenderPostUpdate(body).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.PostUpdateWithBody(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestStorageHandler_PostValueWithBody_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mockserv.NewMockService(ctrl)
	log := mocklog.NewMockLogger(ctrl)
	h := NewStorageHandler(log, svc)

	reqModel := models.GetValueRequest{MType: models.Counter, ID: "hits"}
	reqBody, _ := json.Marshal(reqModel)

	d := int64(42)
	resp := models.GetValueResponse{Delta: &d}

	svc.EXPECT().
		SenderGetValue(reqModel).
		Return(resp, nil)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	h.PostValueWithBody(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var got models.GetValueResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	require.NotNil(t, got.Delta)
	require.Equal(t, int64(42), *got.Delta)
}
