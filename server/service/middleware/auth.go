package middleware

import (
	"net/http"
	"strconv"

	auth "github.com/scott-mescudi/codelet/shared/auth"
)

const ACCESS = 0
const REFRESH = 1

func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		userID, tokenType, err := auth.ValidateHmac(token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if tokenType != ACCESS {
			w.WriteHeader(http.StatusForbidden)
			return			
		}

		r.Header.Add("X-USERID", strconv.Itoa(userID))
		next.ServeHTTP(w, r)
	})
}
