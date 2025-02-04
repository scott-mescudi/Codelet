package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	srv "github.com/scott-mescudi/codelet/service"
	"github.com/scott-mescudi/codelet/service/middleware"
)

func main() {
	if err := os.MkdirAll("/src/logs", 0755); err != nil {
		log.Fatalln("Failed to create logs directory")
	}

	app, clean := srv.NewCodeletServer()
	defer clean()

	server := http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: middleware.CorsMiddleware(app),
	}

	fmt.Println("Server starting on port: ", os.Getenv("APP_PORT"))
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
