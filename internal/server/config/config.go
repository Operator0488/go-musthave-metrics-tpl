package config

import (
	"flag"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/caarlos0/env/v11"
	"strings"
)

func NewServerConfig() (models.ServerConfig, error) {
	return getConfig()
}

func getConfig() (models.ServerConfig, error) {
	conf := models.ServerConfig{}

	port := "localhost:8080"
	filePath := ".storage"
	interval := 0
	restore := false
	dsn := ""

	flag.StringVar(&port, "a", port, "address and port to run server")
	flag.IntVar(&interval, "i", interval, "store interval")
	flag.BoolVar(&restore, "r", restore, "restore data from storage")
	flag.StringVar(&filePath, "f", filePath, "dir for file storage")
	flag.StringVar(&dsn, "d", dsn, "dsn for database")

	_ = env.Parse(&conf)

	if conf.Port != "" {
		port = conf.Port
	}
	if conf.StoreInterval != nil {
		interval = *conf.StoreInterval
	}
	if conf.Restore != nil {
		restore = *conf.Restore
	}
	if conf.FileStoragePath != "" {
		filePath = conf.FileStoragePath
	}
	if conf.DbDsn != "" {
		dsn = conf.DbDsn
	}

	flag.Parse()

	conf.Port = port
	conf.FileStoragePath = strings.TrimSuffix(filePath, "/") + "/wall.txt"
	conf.DbDsn = dsn
	conf.StoreInterval = &interval
	conf.Restore = &restore

	return conf, nil
}
