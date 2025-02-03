package main

import (
	"log"
	"os"

	srv "github.com/scott-mescudi/codelet/service"
)

func main() {
	if err := os.MkdirAll("/src/logs", 0755); err != nil {
		log.Fatalln("Failed to create logs directory")
	}

	srv.NewCodeletServer()
}
