package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
)

func compressGzip(cli *resty.Client, req *resty.Request) error {
	data, _ := json.Marshal(req.Body)

	var wr bytes.Buffer

	wrr := gzip.NewWriter(&wr)
	_, err := wrr.Write(data)
	if err != nil {
		return err
	}
	wrr.Close()

	compressedBody := wr.Bytes()

	req.
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody)

	return nil
}

func loggingRequest(cli *resty.Client, req *resty.Request) error {
	logger.Info("Отправка запроса",
		logger.String("URL", req.URL))

	return nil
}

func loggingResponse(cli *resty.Client, res *resty.Response) error {
	logger.Info("Отправлен запрос",
		logger.Int("status code", res.StatusCode()),
		logger.String("body response", string(res.Body())))

	return nil
}
