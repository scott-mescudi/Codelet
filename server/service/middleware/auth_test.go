package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	auth "github.com/scott-mescudi/codelet/shared/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		setupAuth  func() (string, error)
		expectCode int
		userID     int
	}{
		{
			name: "Valid Token",
			setupAuth: func() (string, error) {
				return auth.GenerateHMac(123, ACCESS, time.Now().Add(2*time.Minute)), nil
			},
			expectCode: http.StatusOK,
			userID:     123,
		},
		{
			name:       "Missing Token",
			token:      "",
			expectCode: http.StatusForbidden,
		},
		{
			name: "Invalid Token",
			setupAuth: func() (string, error) {
				return "invalid.token.string", nil
			},
			expectCode: http.StatusForbidden,
		},
		{
			name: "Wrong Token Type",
			setupAuth: func() (string, error) {
				return auth.GenerateHMac(123, REFRESH, time.Now().Add(2*time.Minute)), nil // Wrong token type
			},
			expectCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			if tt.setupAuth != nil {
				token, _ = tt.setupAuth()
			} else {
				token = tt.token
			}

			req := httptest.NewRequest("GET", "/", nil)
			if token != "" {
				req.Header.Set("Authorization", token)
			}

			rw := httptest.NewRecorder()
			handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			handler.ServeHTTP(rw, req)

			assert.Equal(t, tt.expectCode, rw.Code)

			if tt.expectCode == http.StatusOK {
				assert.Equal(t, strconv.Itoa(tt.userID), req.Header.Get("X-USERID"))
			}
		})
	}
}
