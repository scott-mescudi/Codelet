package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	srv "github.com/scott-mescudi/codelet/service"
	"github.com/scott-mescudi/codelet/service/middleware"
)

func PrintEndpoints(port string) (local string, network string) {
	var localIp string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		if ipNEt, ok := addr.(*net.IPNet); ok && !ipNEt.IP.IsLoopback(){
			if ipNEt.IP.To4() != nil {
				localIp = ipNEt.IP.String()
			}
		}
	}

	return fmt.Sprintf("- Local: http://localhost%s\n", port), fmt.Sprintf("- Local: http://%s%s\n", localIp ,port)
}

func main() {
	if err := os.MkdirAll("/src/logs", 0755); err != nil {
		log.Fatalln("Failed to create logs directory")
	}

	app, clean := srv.NewCodeletServer()
	defer clean()

	port := os.Getenv("APP_PORT")

	local, network := PrintEndpoints(port)
	log.Println(local)
	log.Println(network)

	server := http.Server{
		Addr:    port,
		Handler: middleware.CorsMiddleware(app),
	}

	fmt.Println("Server starting on port: ", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
