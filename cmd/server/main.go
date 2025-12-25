package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/service"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	str := repository.NewMaps(ctx)
	srv := service.NewStorageService(ctx, str)
	mux := handler.NewStorageHandler(ctx, srv)
	rout := handler.NewRoute(mux)

	if err := handler.Listner(":8080", rout); err != nil {
		return err
	}

	return nil
}
