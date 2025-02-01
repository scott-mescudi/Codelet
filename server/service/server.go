package server

import (
	"context"
	"fmt"

	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	snippetMethods "github.com/scott-mescudi/codelet/service/api/snippets"
	userMethods "github.com/scott-mescudi/codelet/service/api/users"
	dataAccess "github.com/scott-mescudi/codelet/service/data_access"
	middleware "github.com/scott-mescudi/codelet/service/middleware"
)

func NewCodeletServer() {
	if err := os.MkdirAll("/src/logs", 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create logs directory")
	}

	file, err := os.OpenFile("/src/logs/codelet_server_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	logger := log.Logger

	app := http.NewServeMux()

	db, err := dataAccess.ConnectToDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		return
	}
	defer db.Close(context.Background())

	query := map[string]string{
		"add_user":                  `INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)`,
		"get_user_password":         `SELECT password_hash, id FROM users WHERE email=$1`,
		"get_user_password_via_id":  `SELECT password_hash FROM users WHERE id=$1`,
		"update_user_password":      `UPDATE users SET password_hash=$1 WHERE id=$2`,
		"add_refresh_token":         `UPDATE users SET refresh_token=$1 WHERE id=$2`,
		"get_refresh_token":         `SELECT refresh_token FROM users WHERE id=$1`,
		"add_snippet":               `INSERT INTO snippets(userid, language, title, code, description, private, tags, created, updated, favorite) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		"get_all_snippet_by_userid": `SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE userid=$1`,
		"get_snippet_by_userid":     `SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE userid=$1 LIMIT $2 OFFSET $3`,
		"get_public_snippet":        `SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE private=false LIMIT $1 OFFSET $2`,
		"delete_snippet":            `DELETE FROM snippets WHERE id=$1`,
	}

	db, err = dataAccess.PrepareStatements(query, db)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to prepared SQL statements")
		return
	}

	logger.Info().Msg("Prepared SQL statements")

	srv := userMethods.UserService{Db: db}
	srv2 := snippetMethods.SnippetService{Db: db, Logger: logger}

	app.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	app.HandleFunc("POST /api/v1/register", srv.Signup)
	app.HandleFunc("POST /api/v1/login", srv.Login)
	app.HandleFunc("GET /api/v1/refresh", srv.Refresh)
	app.Handle("POST /api/v1/update/password", middleware.AuthMiddleware(srv.ChangePassword))
	app.Handle("POST /api/v1/logout", middleware.AuthMiddleware(srv.Logout))

	app.Handle("POST /api/v1/user/snippets", middleware.AuthMiddleware(srv2.AddSnippet))
	app.Handle("DELETE /api/v1/user/snippets/{id}", middleware.AuthMiddleware(srv2.DeleteSnippet))
	app.Handle("GET /api/v1/user/snippets", middleware.AuthMiddleware(srv2.GetUserSnippets))
	app.HandleFunc("GET /api/v1/public/snippets", srv2.GetPublicSnippets)

	logger.Info().Msg(fmt.Sprintf("Server started on port %v", os.Getenv("APP_PORT")))
	if err := http.ListenAndServe(os.Getenv("APP_PORT"), app); err != nil {
		fmt.Println(err)
		return
	}
}
