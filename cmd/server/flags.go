package main

import (
	"flag"
	"fmt"
	"log"
)

func parseFlags(port *string) error {
	flag.StringVar(port, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	log.Println("port server:", port)

	if len(flag.Args()) > 0 {
		return fmt.Errorf("Ошибка, неизвестные флаги: %v", flag.Args())
	}

	return nil
}
