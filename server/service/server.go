package server

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	snippetMethods "github.com/scott-mescudi/codelet/service/api/snippets"
	userMethods "github.com/scott-mescudi/codelet/service/api/users"
	dataAccess "github.com/scott-mescudi/codelet/service/data_access"
	middleware "github.com/scott-mescudi/codelet/service/middleware"
)

func NewCodeletServer() (*http.ServeMux, func()) {
	file, err := os.OpenFile("/src/logs/codelet_server_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	logger := log.Logger

	app := http.NewServeMux()

	logger.Info().Str("Database uri", os.Getenv("DATABASE_URL")).Msg("Trying to connect to database")
	db, err := dataAccess.ConnectToDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		return nil, nil
	}

	clean := func() {
		db.Close()
	}

	logger.Info().Msg("Connected to database")
	srv := userMethods.UserService{Db: db}
	srv2 := snippetMethods.SnippetService{Db: db, Logger: logger}

	app.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	app.HandleFunc("POST /api/v1/register", srv.Signup)
	app.HandleFunc("POST /api/v1/login", srv.Login)
	app.HandleFunc("GET /api/v1/refresh", srv.Refresh)
	app.Handle("GET /api/v1/username", middleware.AuthMiddleware(srv.GetUsernameByID))
	app.Handle("POST /api/v1/update/password", middleware.AuthMiddleware(srv.ChangePassword))
	app.Handle("POST /api/v1/logout", middleware.AuthMiddleware(srv.Logout))
	app.Handle("POST /api/v1/user/snippets", middleware.AuthMiddleware(srv2.AddSnippet))
	app.Handle("DELETE /api/v1/user/snippets/{id}", middleware.AuthMiddleware(srv2.DeleteSnippet))
	app.Handle("GET /api/v1/user/snippets/{id}", middleware.AuthMiddleware(srv2.GetUserSnippetByID))
	app.Handle("GET /api/v1/user/small/snippets", middleware.AuthMiddleware(srv2.GetSmallUserSnippets))
	app.Handle("GET /api/v1/user/snippets", middleware.AuthMiddleware(srv2.GetUserSnippets))
	app.Handle("PUT /api/v1/user/snippets/{id}", middleware.AuthMiddleware(srv2.UpdateUserSnippetByID))
	app.HandleFunc("GET /api/v1/public/snippets", srv2.GetPublicSnippets)

	return app, clean
}
