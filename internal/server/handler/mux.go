package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type StorageHandler struct {
	service service.Service
	log     logger.Logger
}

//go:generate mockgen -source=mux.go -destination=./mocks/mock_srv.go -package=mocks
type Handler interface {
	PostUpdate(http.ResponseWriter, *http.Request)
	PostUpdateWithBody(http.ResponseWriter, *http.Request)
	PostValueWithBody(http.ResponseWriter, *http.Request)
	GetValue(http.ResponseWriter, *http.Request)
	GetValues(http.ResponseWriter, *http.Request)
	returnLogger() logger.Logger
}

func NewRoute(h Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(
		"/update/{type}/{name}/{value}",
		middleConveyor(
			http.HandlerFunc(h.PostUpdate),
			checkPost,
		))

	return mux
}

func NewChiRoute(h Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {

		r.Use(
			logging(h.returnLogger()),
			compressGzip(h.returnLogger()),
			decompressGzip(h.returnLogger()),
		)

		r.Get("/", h.GetValues)

		r.Route("/update", func(r chi.Router) {
			r.Post("/", h.PostUpdateWithBody)
			r.Post("/{type}/{name}/{value}", h.PostUpdate)
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", h.PostValueWithBody)
			r.Get("/{type}/{name}", h.GetValue)
		})
	})

	return r
}

func NewStorageHandler(log logger.Logger, service service.Service) *StorageHandler {
	return &StorageHandler{
		service: service,
		log:     log,
	}
}

func (s *StorageHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	if t != models.Gauge && t != models.Counter {
		errorStatusNotFound(w, fmt.Errorf("ошибка: неправильный тип, %v", t).Error())
		return
	}

	n := chi.URLParam(r, "name")
	if n == "" {
		errorStatusNotFound(w, fmt.Errorf("ошибка: не задано имя").Error())
		return
	}

	req := models.GetValueRequest{
		MType: t,
		ID:    n,
	}

	str, err := s.service.SenderGetValue(req)
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("ошибка: %w", err).Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if str.Value != nil {
		_, err = fmt.Fprint(w, *str.Value)
		if err != nil {
			s.log.Info("Ошибка записи в ResponseWriter",
				zap.Error(err))
		}
	} else if str.Delta != nil {
		_, err = fmt.Fprint(w, *str.Delta)
		if err != nil {
			s.log.Info("Ошибка записи в ResponseWriter",
				zap.Error(err))
		}
	}

}

func (s *StorageHandler) GetValues(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.SenderGetValues()
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("ошибка: %v", err).Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_, err = fmt.Fprint(w, data)
		if err != nil {
			s.log.Info("ошибка записи в ResponseWriter",
				zap.Error(err))
		}
	}
}

func (s *StorageHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	shares := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(shares) != 4 || shares[0] != "update" {
		errorBadRequest(w, fmt.Errorf("ошибка: неправильный запрос, %v", r.URL.Path).Error())
	}

	if shares[1] != models.Gauge && shares[1] != models.Counter {
		errorBadRequest(w, fmt.Errorf("ошибка: неправильный тип, %v", shares[1]).Error())
	}

	req := models.PostUpdateRequest{
		MType:    shares[1],
		ID:       shares[2],
		ValueStr: shares[3],
	}

	err := s.service.SenderPostUpdate(req)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *StorageHandler) PostUpdateWithBody(w http.ResponseWriter, r *http.Request) {
	var data bytes.Buffer
	_, err := data.ReadFrom(r.Body)
	if err != nil {
		s.log.Info("ошибка при чтении body",
			zap.String("err", err.Error()))
		errorBadRequest(w, fmt.Errorf("ошибка: %w", err).Error())
		return
	}

	var req models.PostUpdateRequest
	err = json.Unmarshal(data.Bytes(), &req)
	if err != nil {
		s.log.Info("ошибка при unmarshal",
			zap.String("err", err.Error()))
		errorBadRequest(w, fmt.Errorf("ошибка: %w", err).Error())
		return
	}

	err = s.service.SenderPostUpdate(req)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *StorageHandler) PostValueWithBody(w http.ResponseWriter, r *http.Request) {
	var data bytes.Buffer
	_, err := data.ReadFrom(r.Body)
	if err != nil {
		s.log.Info("ошибка при чтении body",
			zap.String("err", err.Error()),
		)
		errorBadRequest(w, fmt.Errorf("ошибка: %w", err).Error())
		return
	}

	var req models.GetValueRequest
	err = json.Unmarshal(data.Bytes(), &req)
	if err != nil {
		s.log.Info("ошибка при unmarshal",
			zap.String("err", err.Error()),
		)
		errorBadRequest(w, fmt.Errorf("ошибка: %w", err).Error())
		return
	}

	res, err := s.service.SenderGetValue(req)
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("ошибка: %v", err).Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(res)
	if err != nil {
		s.log.Info("ошибка при Marshal",
			zap.Error(err),
		)
		errorInternalServer(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *StorageHandler) returnLogger() logger.Logger {
	return s.log
}
