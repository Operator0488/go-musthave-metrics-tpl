package handler

import (
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type ClientResty struct {
	cli  *resty.Client
	mn   service.Manager
	conf config.AgentConfig
	log  logger.Logger
}

type Sender interface {
	SendRequest()
}

func NewClientResty(log logger.Logger, mn service.Manager, conf config.AgentConfig) *ClientResty {
	cli := resty.New()

	cli.
		OnBeforeRequest(compressGzip).
		OnBeforeRequest(loggingRequest(log)).
		OnAfterResponse(loggingResponse(log))

	return &ClientResty{
		cli:  cli,
		mn:   mn,
		conf: conf,
		log:  log,
	}
}

func (c *ClientResty) SendRequest() {
	m := c.mn.GetMap()

	for k, v := range m {
		req := models.PostUpdateRequest{
			MType: v.Type,
			ID:    k,
		}
		if v.Type == "gauge" {
			node := v.Value
			req.Value = &node
		} else {
			node := int64(v.Value)
			req.Delta = &node
		}

		resp, err := c.cli.R().
			SetHeader("Content-Type", "application/json").
			SetBody(req).
			Post(fmt.Sprintf("http://%s/update/", c.conf.Port))

		if err != nil {
			c.log.Info("Ошибка при отправке запроса",
				zap.Error(err),
				zap.Int("status", resp.StatusCode()))
		}
	}
}
