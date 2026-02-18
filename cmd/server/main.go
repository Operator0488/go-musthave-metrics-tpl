package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/handler"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/service"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	conf, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err = run(ctx, conf); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, conf models.ServerConfig) error {
	l, err := logger.New("info")
	if err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %w", err)
	}

	var str repository.MemStorage

	//conf.DbDsn = "postgres://postgres:yourpasswords@localhost:5432/postgres?sslmode=disable"

	if conf.DbDsn != "" {
		db, err := sql.Open("pgx", conf.DbDsn)
		if err != nil {
			return fmt.Errorf("ошибка подключения к БД: %w", err)
		}
		err = migrations.Migration(db)
		if err != nil {
			l.Info("миграция не прошла",
				zap.Error(err))
		}

		str, err = repository.NewPostgresStorage(
			l.With(zap.String("component", "db")),
			db,
		)
		if err != nil {
			return fmt.Errorf("ошибка инициализации сервера: %w", err)
		}

	} else {
		str, err = repository.NewMaps(
			ctx,
			l.With(zap.String("component", "storage")),
			*conf.Restore,
			conf.FileStoragePath,
			*conf.StoreInterval,
		)

		if err != nil {
			return fmt.Errorf("ошибка инициализации сервера: %w", err)
		}
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
		return fmt.Errorf("ошибка запуска http-сервера: %w", err)
	}
	return nil
}
