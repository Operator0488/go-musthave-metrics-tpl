package handler

import (
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
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
		str = append(str, fmt.Sprintf("http://%s/update/%s/%s/%v", c.conf.Port, v.Type, k, v.Value))
		log.Println(fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%v", v.Type, k, v.Value))
	}

	return str
}

func (c *ClientResty) SendRequest() []error {
	var errors []error

	requests := c.GetRequests()

	for _, req := range requests {
		resp, err := c.cli.R().
			SetHeader("Content-Type", "text/plain").
			Post(req)

		if err != nil || resp.StatusCode() != http.StatusOK {
			errors = append(errors, fmt.Errorf("Ошибка: %v\n Статус ответа: %v ", err, resp.StatusCode()))
		}
	}

	return errors
}
