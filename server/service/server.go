package server

import (
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

	logger.Info().Str("Database uri", os.Getenv("DATABASE_URL")).Msg("Trying to connect to database")
	db, err := dataAccess.ConnectToDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		return
	}
	defer db.Close()
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
