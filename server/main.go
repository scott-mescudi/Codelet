package main

import (


	"github.com/joho/godotenv"
	srv "github.com/scott-mescudi/codelet/service"
)

func main() {
	godotenv.Load(".env")
	srv.NewCodeletServer()
}