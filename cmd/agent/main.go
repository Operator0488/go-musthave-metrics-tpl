package main

import (
	"context"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	conf := config.NewAgentConfig()
	if err := parseFlags(&conf); err != nil {
		log.Fatal(err)
	}

	if err := run(ctx, conf); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, conf config.AgentConfig) error {

	l, err := logger.New("info")
	if err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %w", err)
	}

	mn := service.NewStatsManager()

	client := handler.NewClientResty(
		l.With(zap.String("component", "client")),
		mn,
		conf,
	)

	pollTicker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			client.SendRequest()
		case <-pollTicker.C:
			mn.WriteStats()
		}
	}

}
