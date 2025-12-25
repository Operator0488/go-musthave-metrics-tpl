package handler

import (
	"context"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/service"
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
}

// NewRoute -
func NewRoute(h Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(
		"/update/{type}/{name}/{value}",
		middleConveyor(
			http.HandlerFunc(h.PostUpdate),
			//checkContentType,
			checkPost,
			logging,
		))

	return mux
}

// NewStorageHandler -
func NewStorageHandler(ctx context.Context, service service.Service) *StorageHandler {
	return &StorageHandler{
		service: service,
		ctx:     ctx,
	}
}

// PostUpdate -
func (s *StorageHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	req, err := getUpdateRequest(r.URL.Path)
	if err != nil || req == nil {
		log.Println(err)
		errorBadRequest(w, err.Error())
		return
	}

	err = s.service.Sender(req)
	if err != nil {
		log.Println(err)
		errorBadRequest(w, err.Error())
		return
	}

}

func getUpdateRequest(addr string) (*models.UpdateRequest, error) {
	shares := strings.Split(strings.Trim(addr, "/"), "/")
	if len(shares) != 4 || shares[0] != "update" {
		return nil, fmt.Errorf("Ошибка: неправильный запрос, %v", addr)
	}

	if shares[1] != models.Gauge && shares[1] != models.Counter {
		return nil, fmt.Errorf("Ошибка: неправильный тип, %v", shares[1])
	}

	if shares[2] == "" && shares[3] == "" {
		return nil, fmt.Errorf("Ошибка: неправильное имя или значение, имя: %v, значени: %v", shares[2], shares[3])
	}

	req := models.UpdateRequest{
		Type:  shares[1],
		Name:  shares[2],
		Value: shares[3],
	}

	return &req, nil
}
