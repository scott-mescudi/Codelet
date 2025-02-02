package dataaccess

import (
	"context"
	"fmt"

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

	// Ensure the database is reachable
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

func GetUserPasswordHash(dbConn *pgxpool.Pool, email string) (int, string, error) {
	var id int
	var hash string
	row := dbConn.QueryRow(context.Background(), "SELECT password_hash, id FROM users WHERE email=$1", email)
	if err := row.Scan(&hash, &id); err != nil {
		return -1, "", err
	}

	return id, hash, nil
}

func GetUserPasswordHashViaID(dbConn *pgxpool.Pool, id int) (string, error) {
	var hash string
	row := dbConn.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE id=$1", id)
	if err := row.Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}

func UpdatePassword(dbConn *pgxpool.Pool, passwordHash string, userID int) error {
	_, err := dbConn.Exec(context.Background(), "UPDATE users SET password_hash=$1 WHERE id=$2", passwordHash, userID)
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
