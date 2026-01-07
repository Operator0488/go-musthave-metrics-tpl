package main

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/handler"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"net/http"
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

	if err := http.ListenAndServe(":8080", rout); err != nil {
		return err
	}

	return nil
}
