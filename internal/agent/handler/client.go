package handler

import (
	"context"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	ctx  context.Context
	cli  http.Client
	mn   service.Manager
	conf config.AgentConfig
}

type Sender interface {
	SendRequest() []error
	GetRequests() []string
}

func NewClient(ctx context.Context, mn service.Manager, conf config.AgentConfig) *Client {
	return &Client{
		ctx:  ctx,
		cli:  http.Client{},
		mn:   mn,
		conf: conf,
	}
}

func (c *Client) SendRequest() []error {
	var errors []error

	requests := c.GetRequests()

	for _, req := range requests {
		resp, err := c.cli.Post(req, "text/plain", strings.NewReader(""))
		if err != nil || resp.StatusCode != http.StatusOK {
			errors = append(errors, fmt.Errorf("Ошибка: %v\n Статус ответа: %v ", err, resp.StatusCode))
		}
	}

	return errors
}

func (c *Client) GetRequests() []string {
	m := c.mn.GetMap()
	str := make([]string, 0, len(m))

	for k, v := range m {
		str = append(str, fmt.Sprintf("http://%s/update/%s/%s/%v", c.conf.Port, v.Type, k, v.Value))
		log.Println(fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%v", v.Type, k, v.Value))
	}

	return str
}
