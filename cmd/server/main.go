package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"net/http"
)

func main() {
	ctx := context.Background()

	conf := config.NewServerConfig()
	if err := parseFlags(&conf); err != nil {
		return
	}

	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	if err := run(ctx, conf); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, conf config.ServerConfig) error {
	str := repository.NewMaps()
	srv := service.NewStorageService(str)
	mux := handler.NewStorageHandler(srv)
	rout := handler.NewChiRoute(mux)

	if err := http.ListenAndServe(conf.Port, rout); err != nil {
		return err
	}

	return nil
}
