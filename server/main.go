package main

import (
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	srv "github.com/scott-mescudi/codelet/service"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	f, err := os.Create("/src/logs/cpu.prof")
	if err != nil {
		log.Fatal(err)

	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	srv.NewCodeletServer()
}
