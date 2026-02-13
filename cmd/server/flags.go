package main

import (
	"flag"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/caarlos0/env/v11"
	"strings"
)

func parseFlags(conf *config.ServerConfig) error {
	err := env.Parse(conf)
	if err != nil {
		return fmt.Errorf("ошибка, парсинга env: %w", err)
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

	flag.Parse()

	if len(flag.Args()) > 0 {
		return fmt.Errorf("ошибка, неизвестные флаги: %v", flag.Args())
	}

	return nil
}
