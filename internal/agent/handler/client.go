package handler

import (
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
)

type ClientResty struct {
	cli  *resty.Client
	mn   service.Manager
	conf config.AgentConfig
}

type Sender interface {
	SendRequest() []error
	GetRequests() []string
}

func NewClientResty(mn service.Manager, conf config.AgentConfig) *ClientResty {
	return &ClientResty{
		cli:  resty.New(),
		mn:   mn,
		conf: conf,
	}
}

func (c *ClientResty) GetRequests() []string {
	m := c.mn.GetMap()
	str := make([]string, 0, len(m))

	for k, v := range m {
		str = append(str, fmt.Sprintf("http://%s/update/", c.conf.Port))

		log.Println(fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%v", v.Type, k, v.Value))
	}

	return str
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
			SetBody(&req).
			Post(fmt.Sprintf("http://%s/update/", c.conf.Port))

		logger.Info("Отправка запроса",
			logger.String("URL", resp.Request.URL))

		if err != nil || resp.StatusCode() != http.StatusOK {
			errors = append(errors, fmt.Errorf("Ошибка: %v\n Статус ответа: %v ", err, resp.StatusCode()))
		}

	}

	return errors
}
