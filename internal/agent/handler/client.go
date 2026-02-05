package handler

import (
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"github.com/go-resty/resty/v2"
)

type ClientResty struct {
	cli  *resty.Client
	mn   service.Manager
	conf config.AgentConfig
}

type Sender interface {
	SendRequest() []error
}

func NewClientResty(mn service.Manager, conf config.AgentConfig) *ClientResty {
	cli := resty.New()

	cli.
		OnBeforeRequest(compressGzip).
		OnBeforeRequest(loggingRequest).
		OnAfterResponse(loggingResponse)

	return &ClientResty{
		cli:  cli,
		mn:   mn,
		conf: conf,
	}
}

func (c *ClientResty) SendRequest() []error {
	var errors []error

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
			errors = append(errors, fmt.Errorf("Ошибка при отправке запроса: %v\n ", err, resp.StatusCode()))
		}
	}

	return errors
}
