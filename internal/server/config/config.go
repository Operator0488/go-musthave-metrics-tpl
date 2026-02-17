package config

import (
	"flag"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/caarlos0/env/v11"
	"strings"
)

func NewServerConfig() (models.ServerConfig, error) {
	return getConfig()
}

func getConfig() (models.ServerConfig, error) {
	conf := models.ServerConfig{}

	err := env.Parse(&conf)
	if err != nil {
		return conf, fmt.Errorf("ошибка, парсинга env: %w", err)
	}

	if conf.Port == "" {
		flag.StringVar(&conf.Port, "a", "localhost:8080", "address and port to run server")
	}

	if conf.StoreInterval == nil {
		var si int
		flag.IntVar(&si, "i", 0, "store interval")
		conf.StoreInterval = &si
	}

	if conf.Restore == nil {
		var r bool
		flag.Bool("r", r, "false")
		conf.Restore = &r
	}

	if conf.FileStoragePath == "" {
		flag.StringVar(&conf.FileStoragePath, "f", ".storage", "address and port to run server")
		conf.FileStoragePath = strings.TrimSuffix(conf.FileStoragePath, "/")
		conf.FileStoragePath += "/wall.txt"
	}

	if conf.DbDsn == "" {
		flag.StringVar(&conf.DbDsn, "d", "", "dsn fo database")
	}

	flag.Parse()

	if len(flag.Args()) > 0 {
		return conf, fmt.Errorf("ошибка, неизвестные флаги: %v", flag.Args())
	}

	return conf, nil
}
