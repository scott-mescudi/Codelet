package snippets

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	vsr "github.com/scott-mescudi/codelet/service/api/users"
	dba "github.com/scott-mescudi/codelet/service/data_access"
	"github.com/scott-mescudi/codelet/service/middleware"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(testData string) (*pgxpool.Pool, func(), error) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:latest", postgres.WithDatabase("testdb"), postgres.WithUsername("testAdmin"), postgres.WithPassword("pass1234"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start PostgreSQL container: %v", err)
	}

	time.Sleep(3 * time.Second)

	str, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get PostgreSQL uri: %v", err)
	}

	conn, err := dba.ConnectToDatabase(str)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to PostgreSQL DB: %v", err)
	}

	clean := func() {
		conn.Close()
		pgContainer.Terminate(ctx)
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user', 'moderator')),
		password_hash VARCHAR(255) NOT NULL,
		refresh_token text DEFAULT null,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		clean()
		return nil, nil, fmt.Errorf("failed to create users table: %v", err)
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS snippets (
		id SERIAL PRIMARY KEY,
		userid INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		language VARCHAR(50) NOT NULL,
		favorite boolean DEFAULT false,
		title VARCHAR(255) NOT NULL UNIQUE,
		code BYTEA NOT NULL,
		description TEXT,
		private boolean NOT NULL,
		tags VARCHAR(50)[],
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		clean()
		return nil, nil, fmt.Errorf("failed to create snippets table: %v", err)
	}

	if testData != "" {
		_, err = conn.Exec(ctx, testData)
		if err != nil {
			clean()
			return nil, nil, fmt.Errorf("failed to create test data: %v", err)
		}
	}

	return conn, clean, nil
}

func TestAddUser(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	conn, clean, err := setupTestDB(fmt.Sprintf(`INSERT INTO users (username, email, role, password_hash) VALUES ('fakeuser', 'fakeuser@example.com', 'user', '%s');`, string(hashedPassword)))
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	sp := &vsr.UserService{Db: conn, Logger: zerolog.New(os.Stdout)}

	body, err := json.Marshal(vsr.UserLogin{Email: "fakeuser@example.com", Password: "hashedpassword123"})
	if err != nil {
		t.Fatal(err)
	}

	loginReq := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()

	sp.Login(loginRec, loginReq)

	var rr struct {
		Token string `json:"access_token"`
	}

	if err := json.NewDecoder(loginRec.Body).Decode(&rr); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		info     Snippet
		expected int
	}{
		{
			name: "Valid request",
			info: Snippet{
				Language:    "go",
				Title:       "go test",
				Code:        "fmt.Println('hello)",
				Favorite:    false,
				Private:     false,
				Tags:        []string{"sigma", "wobc"},
				Description: "wljkhf",
			},
			expected: http.StatusCreated,
		},
		{
			name: "Empty title request",
			info: Snippet{
				Language:    "go",
				Title:       "",
				Code:        "fmt.Println('hello)",
				Favorite:    false,
				Private:     false,
				Tags:        []string{"sigma", "wobc"},
				Description: "wljkhf",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Empty language request",
			info: Snippet{
				Language:    "",
				Title:       "go test",
				Code:        "fmt.Println('hello)",
				Favorite:    false,
				Private:     false,
				Tags:        []string{"sigma", "wobc"},
				Description: "wljkhf",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Empty code request",
			info: Snippet{
				Language:    "go",
				Title:       "go test",
				Code:        "",
				Favorite:    false,
				Private:     false,
				Tags:        []string{"sigma", "wobc"},
				Description: "wljkhf",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Too large code request",
			info: Snippet{
				Language:    "go",
				Title:       "go test",
				Code:        string(make([]byte, 5000)),
				Favorite:    false,
				Private:     false,
				Tags:        []string{"sigma", "wobc"},
				Description: "wljkhf",
			},
			expected: http.StatusRequestEntityTooLarge,
		},
	}

	app := SnippetService{Db: conn, Logger: zerolog.New(os.Stdout)}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			body, err := json.Marshal(tt.info)
			if err != nil {
				t.Error(err)
			}
			req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", rr.Token)

			handler := middleware.AuthMiddleware(http.HandlerFunc(app.AddSnippet))
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expected {
				t.Errorf("expected status code %d, got %d", tt.expected, rec.Code)
			}
		})
	}

	t.Run("Invalid json payload", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body, err := json.Marshal("")
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.AddSnippet))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("Wrong content type", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body, err := json.Marshal("")
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "json/text")
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.AddSnippet))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Missing auth token", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body, err := json.Marshal("")
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		app.AddSnippet(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Wrong X-USERID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body, err := json.Marshal("")
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-USERID", "")

		app.AddSnippet(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Missing X-USERID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body, err := json.Marshal("")
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		app.AddSnippet(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetUserSnippets(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	conn, clean, err := setupTestDB(fmt.Sprintf(`INSERT INTO users (username, email, role, password_hash) VALUES ('fakeuser', 'fakeuser@example.com', 'user', '%s');`, string(hashedPassword)))
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	sp := &vsr.UserService{Db: conn, Logger: zerolog.New(os.Stdout)}
	if err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(vsr.UserLogin{Email: "fakeuser@example.com", Password: "hashedpassword123"})
	if err != nil {
		t.Fatal(err)
	}

	loginReq := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()

	sp.Login(loginRec, loginReq)

	var rr struct {
		Token string `json:"access_token"`
	}

	if err := json.NewDecoder(loginRec.Body).Decode(&rr); err != nil {
		t.Fatal(err)
	}

	app := SnippetService{Db: conn, Logger: zerolog.New(os.Stdout)}

	info := Snippet{
		Language:    "go",
		Title:       "go test",
		Code:        "fmt.Println('hello)",
		Favorite:    false,
		Private:     false,
		Tags:        []string{"sigma", "wobc"},
		Description: "wljkhf",
	}

	rec := httptest.NewRecorder()
	body, err = json.Marshal(info)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("POST", "/api/v1/user/snippets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", rr.Token)

	handler := middleware.AuthMiddleware(http.HandlerFunc(app.AddSnippet))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	t.Run("Valid Request", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets?page=1&limit=5", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.GetUserSnippets))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Missing params", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.GetUserSnippets))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets?page=3&limit=5", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.GetUserSnippets))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("limit too high", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets?page=1&limit=2048", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.GetUserSnippets))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("No userid header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets?page=-2&limit=-2048", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		app.GetUserSnippets(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("limit and page too low", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/user/snippets?page=-2&limit=-2048", bytes.NewBuffer(body))
		req.Header.Set("Authorization", rr.Token)

		handler := middleware.AuthMiddleware(http.HandlerFunc(app.GetUserSnippets))
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
