package main

import (
	"flag"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/config"
	"github.com/caarlos0/env/v6"
)

func parseFlags(conf *config.ServerConfig) error {

	_ = env.Parse(conf)

	if conf.Port == "" {
		flag.StringVar(&conf.Port, "a", "localhost:8080", "address and port to run server")
		flag.Parse()
	}

	if len(flag.Args()) > 0 {
		return fmt.Errorf("Ошибка, неизвестные флаги: %v", flag.Args())
	}

	return nil
}
