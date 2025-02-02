package users

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
	dba "github.com/scott-mescudi/codelet/service/data_access"
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

func TestSignup(t *testing.T) {
	conn, clean, err := setupTestDB(`INSERT INTO users (username, email, role, password_hash) VALUES ('fakeuser', 'fakeuser@example.com', 'user', 'hashedpassword123');`)
	if err != nil {
		t.Error(err)
		return
	}
	defer clean()

	app := &UserService{Db: conn, Logger: zerolog.New(os.Stdout)}

	validTests := []struct {
		name     string
		user     UserSignup
		expected int
	}{
		{
			name: "Valid user",
			user: UserSignup{
				Username: "jacky",
				Email:    "sigma@sigma.com",
				Password: "flstudiosucks",
				Role:     "admin",
			},

			expected: http.StatusCreated,
		},
		{
			name: "Invalid email",
			user: UserSignup{
				Username: "j",
				Email:    "sigma@.com",
				Password: "flstudiosucks",
				Role:     "admin",
			},

			expected: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			user: UserSignup{
				Username: "jacky",
				Email:    "sigma@sigma.com",
				Password: "",
				Role:     "admin",
			},

			expected: http.StatusBadRequest,
		},
		{
			name: "Missing username",
			user: UserSignup{
				Username: "",
				Email:    "sigma@sigma.com",
				Password: "ksghd",
				Role:     "user",
			},

			expected: http.StatusBadRequest,
		},
		{
			name: "Missing role",
			user: UserSignup{
				Username: "dsnf",
				Email:    "sigma@sigma.com",
				Password: "skdj",
				Role:     "",
			},

			expected: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			user: UserSignup{
				Username: "fakeuser",
				Email:    "fakeuser@example.com",
				Password: "hashedpassword123",
				Role:     "admin",
			},

			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			app.Signup(rec, req)

			if rec.Code != tt.expected {
				t.Errorf("Expected status %v, got %v", tt.expected, rec.Code)
			}
		})
	}

	t.Run("Malformed json", func(t *testing.T) {
		body, err := json.Marshal("")
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Signup(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %v, got %v", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("Invalid Content-Type", func(t *testing.T) {
		body, err := json.Marshal("")
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		app.Signup(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %v, got %v", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err)
	}

	conn, clean, err := setupTestDB(fmt.Sprintf(`INSERT INTO users (username, email, role, password_hash) VALUES ('fakeuser', 'fakeuser@example.com', 'user', '%s');`, string(hashedPassword)))
	if err != nil {
		t.Error(err)
		return
	}
	defer clean()

	app := &UserService{Db: conn, Logger: zerolog.New(os.Stdout)}

	validTests := []struct {
		name     string
		info     UserLogin
		expected int
	}{
		{
			name:     "Valid login",
			info:     UserLogin{Email: "fakeuser@example.com", Password: "hashedpassword123"},
			expected: http.StatusOK,
		},
		{
			name:     "Invalid email",
			info:     UserLogin{Email: "fake1user@example.com", Password: "hashedpassword123"},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Invalid password",
			info:     UserLogin{Email: "fake1user@example.com", Password: "hashedpssword123"},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Missing email",
			info:     UserLogin{Email: "", Password: "hashedpassword123"},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Missing password",
			info:     UserLogin{Email: "fake1user@example.com", Password: ""},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.info)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			app.Login(rec, req)

			if rec.Code != tt.expected {
				t.Errorf("Expected status %v, got %v", tt.expected, rec.Code)
			}
		})
	}

	t.Run("Malformed json", func(t *testing.T) {
		body, err := json.Marshal("")
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Login(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %v, got %v", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("Invalid Content-Type", func(t *testing.T) {
		body, err := json.Marshal("")
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		app.Login(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %v, got %v", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

// add more tests like expired cookie, wrong type, invalid cookie
func TestRefresh(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("hashedpassword123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	conn, clean, err := setupTestDB(fmt.Sprintf(`INSERT INTO users (username, email, role, password_hash) VALUES ('fakeuser', 'fakeuser@example.com', 'user', '%s');`, string(hashedPassword)))
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	app := &UserService{Db: conn, Logger: zerolog.New(os.Stdout)}

	body, err := json.Marshal(UserLogin{Email: "fakeuser@example.com", Password: "hashedpassword123"})
	if err != nil {
		t.Fatal(err)
	}

	loginReq := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()

	app.Login(loginRec, loginReq)

	var info struct {
		Token string `json:"access_token"`
	}
	if err := json.NewDecoder(loginRec.Body).Decode(&info); err != nil {
		t.Fatal(err)
	}

	t.Run("Valid Refresh", func(t *testing.T) {
		cookies := loginRec.Result().Cookies()
		if len(cookies) == 0 {
			t.Fatal("No cookies set in login response")
		}

		refreshReq := httptest.NewRequest("POST", "/api/v1/refresh", nil)
		refreshReq.Header.Set("Authorization", info.Token)
		for _, cookie := range cookies {
			refreshReq.AddCookie(cookie)
		}
		refreshRec := httptest.NewRecorder()

		app.Refresh(refreshRec, refreshReq)
		if refreshRec.Code != http.StatusOK {
			t.Errorf("Expected 200 got %v", refreshRec.Code)
		}
	})
}
