package middleware

import (
	"net/http"
	"os"
)

var CORS_ORIGIN = os.Getenv("CORS_ORIGIN")

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", CORS_ORIGIN) // Replace with your frontend URL, in testing with postman set to *
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials (cookies, etc.)

		// If the request method is OPTIONS (preflight), just return No Content
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the actual request
		next.ServeHTTP(w, r)
	})
}
