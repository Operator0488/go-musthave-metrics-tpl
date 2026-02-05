package handler

import (
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type StorageHandler struct {
	service service.Service
}

type Handler interface {
	PostUpdate(http.ResponseWriter, *http.Request)
	PostUpdateWithBody(http.ResponseWriter, *http.Request)
	PostValueWithBody(http.ResponseWriter, *http.Request)
	GetValue(http.ResponseWriter, *http.Request)
	GetValues(http.ResponseWriter, *http.Request)
}

// NewRoute -
func NewRoute(h Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(
		"/update/{type}/{name}/{value}",
		middleConveyor(
			http.HandlerFunc(h.PostUpdate),
			//checkPost,
			//logging,
		))

	return mux
}

// NewChiRoute -
func NewChiRoute(h Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {

		r.Use(
			logging,
			compressGzip,
			decompressGzip)

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

func NewStorageHandler(service service.Service) *StorageHandler {
	return &StorageHandler{
		service: service,
	}
}

func (s *StorageHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	if t != models.Gauge && t != models.Counter {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: неправильный тип, %v", t).Error())
		return
	}

	n := chi.URLParam(r, "name")
	if n == "" {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: не задано имя, %v", n).Error())
		return
	}

	req := models.GetValueRequest{
		MType: t,
		ID:    n,
	}

	str, err := s.service.SenderGetValue(req)
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	if str.Value != nil {
		writeText(w, *str.Value)
	}

	if str.Delta != nil {
		writeText(w, *str.Delta)
	}

}

func (s *StorageHandler) GetValues(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.SenderGetValues()
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	writeText(w, data)
}

func (s *StorageHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	req, err := postUpdateRequestAddr(r.URL.Path)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	err = s.service.SenderPostUpdate(req)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	setHeader200(w)
}

func (s *StorageHandler) PostUpdateWithBody(w http.ResponseWriter, r *http.Request) {
	req, err := getUpdateRequestBody(r.Body)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	err = s.service.SenderPostUpdate(req)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	setHeader200(w)
}

func (s *StorageHandler) PostValueWithBody(w http.ResponseWriter, r *http.Request) {
	req, err := getValueRequestBody(r.Body)
	if err != nil {
		errorBadRequest(w, err.Error())
		return
	}

	res, err := s.service.SenderGetValue(req)
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	writeJson(w, res)
}
