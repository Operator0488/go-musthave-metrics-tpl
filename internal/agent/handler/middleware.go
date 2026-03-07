package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func compressGzip(cli *resty.Client, req *resty.Request) error {
	data, err := json.Marshal(req.Body)
	if err != nil {
		return fmt.Errorf("ошибка json.Marshal: %w", err)
	}

	var wr bytes.Buffer

	wrr := gzip.NewWriter(&wr)
	_, err = wrr.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка gzip.NewWriter.Write: %w", err)
	}
	err = wrr.Close()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии writer: %w", err)
	}

	compressedBody := wr.Bytes()

	req.
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody)

	return nil
}

func loggingRequest(log logger.Logger) func(*resty.Client, *resty.Request) error {
	return func(cli *resty.Client, req *resty.Request) error {
		log.Info("Отправка запроса",
			zap.String("url", req.URL),
			zap.String("method", req.Method),
		)
		return nil
	}
}
func loggingResponse(log logger.Logger) func(*resty.Client, *resty.Response) error {
	return func(cli *resty.Client, res *resty.Response) error {
		log.Info("Получен ответ",
			zap.Int("status", res.StatusCode()),
			zap.ByteString("body", res.Body()),
		)
		return nil
	}
}
