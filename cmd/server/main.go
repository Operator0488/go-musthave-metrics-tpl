package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	conf := config.NewServerConfig()
	if err := parseFlags(&conf); err != nil {
		log.Fatal(err)
	}

	if err := run(ctx, conf); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, conf config.ServerConfig) error {
	l, err := logger.New("info")
	if err != nil {
		return err
	}

	str, err := repository.NewMaps(
		l.With(zap.String("component", "storage")),
		*conf.Restore,
		conf.FileStoragePath,
		*conf.StoreInterval,
	)
	if err != nil {
		return err
	}

	srv := service.NewStorageService(
		l.With(zap.String("component", "service")),
		str,
	)

	mux := handler.NewStorageHandler(
		l.With(zap.String("component", "handler")),
		srv,
	)

	rout := handler.NewChiRoute(mux)

	if err = http.ListenAndServe(conf.Port, rout); err != nil {
		return err
	}

	return nil
}
