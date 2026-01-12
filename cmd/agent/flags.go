package main

import (
	"flag"
	"fmt"
	"log"
)

func parseFlags(port *string, repInterval *int, pollInterval *int) error {

	flag.StringVar(port, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(repInterval, "r", 10, "report interval")
	flag.IntVar(pollInterval, "p", 2, "poll interval")
	flag.Parse()

	log.Println("port:", port)
	if len(flag.Args()) > 0 {
		return fmt.Errorf("Ошибка, неизвестные флаги: %v", flag.Args())
	}

	return nil
}
