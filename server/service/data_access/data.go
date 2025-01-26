package dataaccess

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ConnectToDatabase(postgresURI string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), postgresURI)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func PrepareStatements(queries map[string]string, conn *pgx.Conn) (*pgx.Conn, error) {
	for name, query := range queries {
		_, err := conn.Prepare(context.Background(), name, query)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func AddUser(dbConn *pgx.Conn, username, email, role, password string) error {
	_, err := dbConn.Exec(context.Background(), "add_user", username, email, role, password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserPasswordHash(dbConn *pgx.Conn, email string) (int, string, error) {
	var id int
	var hash string
	row := dbConn.QueryRow(context.Background(), "get_user_password", email)
	if err := row.Scan(&hash, &id); err != nil {
		return -1, "", err
	}

	return id, hash, nil
}

func GetUserPasswordHashViaID(dbConn *pgx.Conn, id int) (string, error) {
	var hash string
	row := dbConn.QueryRow(context.Background(), "get_user_password_via_id", id)
	if err := row.Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}


func UpdatePassword(dbConn *pgx.Conn, passwordHash string, userID int) error {
	_, err := dbConn.Exec(context.Background(), "update_user_password", passwordHash, userID)
	if err != nil {
		return err
	}

	return nil
}
