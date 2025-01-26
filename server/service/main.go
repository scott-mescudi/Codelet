package main

import (
	"context"
	"log"
	"net/http"

	userMethods "github.com/scott-mescudi/codelet/service/api/users"
	snippetMethods "github.com/scott-mescudi/codelet/service/api/snippets"
	dataAccess "github.com/scott-mescudi/codelet/service/data_access"
	middleware "github.com/scott-mescudi/codelet/service/middleware"
)



func main() {
	app := http.NewServeMux()

	db, err := dataAccess.ConnectToDatabase("postgresql://admin:password123@localhost:3100/codelet_database")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close(context.Background())

	query := map[string]string{
		"add_user":             `INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)`,
		"get_user_password":    `SELECT password_hash, id FROM users WHERE email=$1`,
		"get_user_password_via_id":    `SELECT password_hash FROM users WHERE id=$1`,
		"update_user_password": `UPDATE users SET password_hash=$1 WHERE id=$2`,
		"add_refresh_token" : `UPDATE users SET refresh_token=$1 WHERE id=$2`,
		"get_refresh_token" : `SELECT refresh_token FROM users WHERE id=$1`,
		"add_snippet" : `INSERT INTO snippets(userid, language, title, code, description, private, tags, created, updated) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		"get_all_snippet_by_userid" : `SELECT id, language, title, code, description, private, tags, created, updated FROM snippets WHERE userid=$1`,
		"get_snippet_by_userid" : `SELECT id, language, title, code, description, private, tags, created, updated FROM snippets WHERE userid=$1 LIMIT $2 OFFSET $3`,
	}

	db, err = dataAccess.PrepareStatements(query, db)
	if err != nil {
		log.Fatalln(err)
	}

	srv := userMethods.UserService{Db: db}
	srv2 := snippetMethods.SnippetService{Db: db}

	app.HandleFunc("POST /api/v1/register", srv.Signup)
	app.HandleFunc("POST /api/v1/login", srv.Login)
	app.HandleFunc("GET /api/v1/refresh", srv.Refresh)
	app.Handle("POST /api/v1/update/password", middleware.AuthMiddleware(srv.ChangePassword)) 
	app.Handle("POST /api/v1/logout", middleware.AuthMiddleware(srv.Logout))
	app.Handle("POST /api/v1/snippets", middleware.AuthMiddleware(srv2.AddSnippet))
	app.Handle("GET /api/v1/snippets", middleware.AuthMiddleware(srv2.GetAllSnippetsFromUser))      

	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalln(err)
	}

}
