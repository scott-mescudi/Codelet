package dataaccess

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase(postgresURI string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(postgresURI)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return dbPool, nil
}

func AddUser(dbConn *pgxpool.Pool, username, email, role, password string) error {
	_, err := dbConn.Exec(context.Background(), "INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)", username, email, role, password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserPasswordHashAndLastLogin(dbConn *pgxpool.Pool, email string) (int, string, *time.Time, error) {
	var id int
	var hash string
	var ll pgtype.Timestamptz
	row := dbConn.QueryRow(context.Background(), "SELECT password_hash, last_login, id FROM users WHERE email=$1", email)
	if err := row.Scan(&hash, &ll, &id); err != nil {
		return -1, "", nil, err
	}

	var lastLogin *time.Time
	if ll.Valid {
		lastLogin = &ll.Time
	}

	return id, hash, lastLogin, nil
}

func GetUserPasswordHashViaID(dbConn *pgxpool.Pool, id int) (string, error) {
	var hash string
	row := dbConn.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE id=$1", id)
	if err := row.Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}

func UpdatePassword(dbConn *pgxpool.Pool, passwordHash string, updatedAt time.Time, userID int) error {
	_, err := dbConn.Exec(context.Background(), "UPDATE users SET password_hash=$1, updated=$2 WHERE id=$3", passwordHash, updatedAt, userID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTokenAndLoginTime(dbConn *pgxpool.Pool, acessToken string, loginTime time.Time, userID int) error {
	_, err := dbConn.Exec(context.Background(), "UPDATE users SET refresh_token=$1, last_login=$2 WHERE id=$3", acessToken, loginTime, userID)
	if err != nil {
		return err
	}

	return nil
}

func AddRefreshToken(dbConn *pgxpool.Pool, acessToken string, userID int) error {
	_, err := dbConn.Exec(context.Background(), "UPDATE users SET refresh_token=$1 WHERE id=$2", acessToken, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetRefreshToken(dbConn *pgxpool.Pool, userID int) (string, error) {
	var refreshToken string
	row := dbConn.QueryRow(context.Background(), "SELECT refresh_token FROM users WHERE id=$1", userID)
	if err := row.Scan(&refreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}
