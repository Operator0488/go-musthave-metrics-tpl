package handler

import (
	"compress/gzip"
	"encoding/json"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service/mocks"
	mocklog "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger/mocks"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestClientResty_SendRequest_SendsGaugeAndCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	type gotReq struct {
		Header http.Header
		Body   models.PostUpdateRequest
	}

	var (
		mu   sync.Mutex
		got  []gotReq
		mapa = make(map[string]models.PostUpdateRequest)
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/update/", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// тело может быть gzip
		var b []byte
		var err error

		if r.Header.Get("Content-Encoding") == "gzip" {
			zr, err := gzip.NewReader(r.Body)
			require.NoError(t, err)
			defer zr.Close()

			b, err = io.ReadAll(zr)
			require.NoError(t, err)
		} else {
			b, err = io.ReadAll(r.Body)
			require.NoError(t, err)
		}

		var req models.PostUpdateRequest
		require.NoError(t, json.Unmarshal(b, &req))

		mu.Lock()
		got = append(got, gotReq{Header: r.Header.Clone(), Body: req})
		mapa[req.ID] = req
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	hostPort := strings.TrimPrefix(srv.URL, "http://")

	conf := config.AgentConfig{Port: hostPort}

	mockmng := mocks.NewMockManager(ctrl)

	mockmng.EXPECT().GetMap().Return(map[string]*model.Stat{
		"lastgc": {Type: "gauge", Value: 12.5},
		"count":  {Type: "counter", Value: 7}, // Value float64 -> Delta int64(7)
	})

	c := NewClientResty(log, mockmng, conf)

	c.SendRequest()

	reqCPU, ok := mapa["lastgc"]
	require.True(t, ok)
	require.Equal(t, "gauge", reqCPU.MType)
	require.Equal(t, "lastgc", reqCPU.ID)
	require.NotNil(t, reqCPU.Value)
	require.Equal(t, 12.5, *reqCPU.Value)
	require.Nil(t, reqCPU.Delta)

	reqHits, ok := mapa["count"]
	require.True(t, ok)
	require.Equal(t, "counter", reqHits.MType)
	require.Equal(t, "count", reqHits.ID)
	require.NotNil(t, reqHits.Delta)
	require.Equal(t, int64(7), *reqHits.Delta)
	require.Nil(t, reqHits.Value)
}
