package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service"
	"time"
)

const (
	pollInterval   = 2 //секунды
	reportInterval = 10
)

func main() {
	ctx := context.Background()
	if err := run(ctx); len(err) != 0 {
		panic(err)
	}
}

func run(ctx context.Context) []error {

	mn := service.NewStatsManager(ctx)
	client := handler.NewClient(ctx, mn)

	pollTicker := time.NewTicker(pollInterval * time.Second)
	reportTicker := time.NewTicker(reportInterval * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			err := client.SendRequest()
			if len(err) != 0 {
				return err
			}

		case <-pollTicker.C:
			mn.WriteStats()

		}
	}

}
