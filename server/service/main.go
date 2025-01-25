package main

import (
	"context"
	"log"
	"net/http"

	methods "github.com/scott-mescudi/codelet/service/api"
	dataAccess "github.com/scott-mescudi/codelet/service/data_access"
)

func main() {
	app := http.NewServeMux()

	app.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("im alive"))
	})

	db, err := dataAccess.ConnectToDatabase("postgresql://admin:password123@localhost:3100/codelet_database")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close(context.Background())

	query := map[string]string{
		"add_user": `INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)`,
	}

	db, err = dataAccess.PrepareStatements(query, db)
	if err != nil {
		log.Fatalln(err)
	}

	srv := methods.Server{Db: db}

	app.HandleFunc("/api/v1/register", srv.Signup)
	app.HandleFunc("/api/v1/login", srv.Login)
	app.HandleFunc("/api/v1/update/password", srv.ChangePassword)
	
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalln(err)
	}
	
}