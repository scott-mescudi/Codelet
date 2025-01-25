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
