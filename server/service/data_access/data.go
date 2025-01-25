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


func  PrepareStatements(queries map[string]string, conn *pgx.Conn) (*pgx.Conn, error) {
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
		return  err
	}

	return nil
}
