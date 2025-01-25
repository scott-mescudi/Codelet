package main

import (
	"log"
	"net/http"
)

func main() {
	app := http.NewServeMux()

	app.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("im alive"))
	})
	
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalln(err)
	}
	
}