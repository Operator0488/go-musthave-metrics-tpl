package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
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

		r.Use(logging)
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

// NewStorageHandler -
func NewStorageHandler(service service.Service) *StorageHandler {
	return &StorageHandler{
		service: service,
	}
}

// GetValue -
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

// GetValues -
func (s *StorageHandler) GetValues(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.SenderGetValues()
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	writeText(w, data)
}

// PostUpdate -
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
}

func postUpdateRequestAddr(addr string) (models.PostUpdateRequest, error) {
	shares := strings.Split(strings.Trim(addr, "/"), "/")
	if len(shares) != 4 || shares[0] != "update" {
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: неправильный запрос, %v", addr)
	}

	if shares[1] != models.Gauge && shares[1] != models.Counter {
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: неправильный тип, %v", shares[1])
	}

	req := models.PostUpdateRequest{
		MType:    shares[1],
		ID:       shares[2],
		ValueStr: shares[3],
	}

	return req, nil
}

func writeText(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_, er := fmt.Fprint(w, data)
		if er != nil {
			panic(er)
		}
	}
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

}

func getUpdateRequestBody(body io.ReadCloser) (models.PostUpdateRequest, error) {
	var data bytes.Buffer
	_, err := data.ReadFrom(body)
	if err != nil {
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	var req models.PostUpdateRequest

	err = json.Unmarshal(data.Bytes(), &req)
	if err != nil {
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	return req, nil
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

func getValueRequestBody(body io.ReadCloser) (models.GetValueRequest, error) {
	var data bytes.Buffer
	_, err := data.ReadFrom(body)
	if err != nil {
		logger.Info("ошибка при чтении body",
			logger.String("err", err.Error()))
		return models.GetValueRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	var req models.GetValueRequest

	err = json.Unmarshal(data.Bytes(), &req)
	if err != nil {
		logger.Info("ошибка при анмаршл",
			logger.String("err", err.Error()))
		return models.GetValueRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	return req, nil
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		b, er := json.Marshal(data)
		if er != nil {
			panic(er)
		}
		_, er = w.Write(b)
		if er != nil {
			panic(er)
		}
	}
}
