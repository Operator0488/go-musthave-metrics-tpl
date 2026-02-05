package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	conf := config.NewAgentConfig()
	if err := parseFlags(&conf); err != nil {
		log.Println(err)
		return
	}

	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	if err := run(ctx, conf); len(err) != 0 {
		panic(err)
	}
}

func run(ctx context.Context, conf config.AgentConfig) []error {

	mn := service.NewStatsManager()
	client := handler.NewClientResty(mn, conf)

	pollTicker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			err := client.SendRequest()
			if len(err) != 0 {
				for _, er := range err {
					logger.Info("Ошибка при отправке запроса",
						logger.String("err", er.Error()))
				}
			}

		case <-pollTicker.C:
			mn.WriteStats()
		}
	}

}
