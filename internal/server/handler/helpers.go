package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"io"
	"net/http"
	"strings"
)

func setHeader200(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
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

func writeText(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_, er := fmt.Fprint(w, data)
		if er != nil {
			panic(er)
		}
	}
}

func getUpdateRequestBody(body io.ReadCloser) (models.PostUpdateRequest, error) {
	var data bytes.Buffer
	_, err := data.ReadFrom(body)
	if err != nil {
		logger.Info("ошибка getUpdateRequestBody",
			logger.String("err", err.Error()))
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	var req models.PostUpdateRequest

	err = json.Unmarshal(data.Bytes(), &req)
	if err != nil {
		logger.Info("ошибка getUpdateRequestBody",
			logger.String("err", err.Error()))
		return models.PostUpdateRequest{}, fmt.Errorf("Ошибка: %v", err)
	}

	return req, nil
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
