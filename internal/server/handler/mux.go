package handler

import (
	"context"
	"encoding/json"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

type StorageHandler struct {
	ctx     context.Context
	service service.Service
}

type Handler interface {
	PostUpdate(http.ResponseWriter, *http.Request)
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
			checkPost,
			logging,
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
			r.Post("/{type}/{name}/{value}", h.PostUpdate)
		})

		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", h.GetValue)
		})

	})

	return r
}

// NewStorageHandler -
func NewStorageHandler(ctx context.Context, service service.Service) *StorageHandler {
	return &StorageHandler{
		service: service,
		ctx:     ctx,
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
		Type: t,
		Name: n,
	}

	str, err := s.service.SenderGetValue(&req)
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	write(w, str)
}

// GetValues -
func (s *StorageHandler) GetValues(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.SenderGetValues()
	if err != nil {
		errorStatusNotFound(w, fmt.Errorf("Ошибка: %v", err).Error())
		return
	}

	write(w, data)
}

// PostUpdate -
func (s *StorageHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	req, err := postUpdateRequest(r.URL.Path)
	if err != nil || req == nil {
		log.Println(err)
		errorBadRequest(w, err.Error())
		return
	}

	err = s.service.SenderPostUpdate(req)
	if err != nil {
		log.Println(err)
		errorBadRequest(w, err.Error())
		return
	}

}

func postUpdateRequest(addr string) (*models.PostUpdateRequest, error) {
	shares := strings.Split(strings.Trim(addr, "/"), "/")
	if len(shares) != 4 || shares[0] != "update" {
		return nil, fmt.Errorf("Ошибка: неправильный запрос, %v", addr)
	}

	if shares[1] != models.Gauge && shares[1] != models.Counter {
		return nil, fmt.Errorf("Ошибка: неправильный тип, %v", shares[1])
	}

	req := models.PostUpdateRequest{
		Type:  shares[1],
		Name:  shares[2],
		Value: shares[3],
	}

	return &req, nil
}

func write(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_, er := fmt.Fprint(w, data)
		if er != nil {
			panic(er)
		}
	}
}

func writeData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		er := json.NewEncoder(w).Encode(data)
		if er != nil {
			panic(er)
		}
	}
}
