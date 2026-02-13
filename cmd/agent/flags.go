package main

import (
	"flag"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/caarlos0/env/v6"
)

func parseFlags(agent *config.AgentConfig) error {
	_ = env.Parse(agent)

	if agent.Port == "" {
		flag.StringVar(&agent.Port, "a", "localhost:8080", "address and port to run server")
	}
	if agent.PollInterval == 0 {
		flag.IntVar(&agent.PollInterval, "p", 2, "report interval")
	}
	if agent.ReportInterval == 0 {
		flag.IntVar(&agent.ReportInterval, "r", 10, "report interval")
	}
	flag.Parse()

	if len(flag.Args()) > 0 {
		return fmt.Errorf("Ошибка, неизвестные флаги: %v", flag.Args())
	}

	return nil
}
