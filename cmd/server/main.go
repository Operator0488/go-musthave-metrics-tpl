package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	conf := config.NewServerConfig()
	if err := parseFlags(&conf); err != nil {
		panic(err)
	}

	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	if err := run(ctx, conf); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, conf config.ServerConfig) error {
	str, err := repository.NewMaps(*conf.Restore, conf.FileStoragePath, *conf.StoreInterval)
	if err != nil {
		return err
	}
	srv := service.NewStorageService(str, *conf.StoreInterval)
	mux := handler.NewStorageHandler(srv)
	rout := handler.NewChiRoute(mux)

	go func() {
		if *conf.StoreInterval > 0 {
			tick := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
			for {
				select {
				case <-tick.C:
					err = str.Snapshot()
					logger.Info("Не удалось сделать снэп",
						logger.Error(err))
				}
			}

		}
	}()

	if err = http.ListenAndServe(conf.Port, rout); err != nil {
		return err
	}

	return nil
}
