package dataaccess

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	cmp "github.com/scott-mescudi/codelet/shared/compression"
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

func AddRefreshToken(dbConn *pgx.Conn, acessToken string, userID int) error {
	_, err := dbConn.Exec(context.Background(), "add_refresh_token", acessToken, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetRefreshToken(dbConn *pgx.Conn, userID int) (string, error) {
	var refreshToken string
	row := dbConn.QueryRow(context.Background(), "get_refresh_token", userID)
	if err := row.Scan(&refreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func AddSnippet(dbConn *pgx.Conn, userID int, language, description, title string, code string, private bool, tags []string, created time.Time, updated time.Time) error {
	compressed, err := cmp.CompressZSTD([]byte(code))
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(context.Background(), "add_snippet", userID, language, title, compressed, description, private, tags, created, updated)
	if err != nil {
		return err
	}

	return nil
}

func GetSnippetsByUserID(dbConn *pgx.Conn, userID, limit, offset int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "get_snippet_by_userid", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet DBsnippet
		var code []byte
		err := rows.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description, &snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated)
		if err != nil {
			return nil, err
		}

		decompressedData, err := cmp.DecompressZSTD(code)
		if err != nil {
			return nil, err
		}

		snippet.Code = string(decompressedData)
		data = append(data, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil

}

func GetAllSnippetsByUserID(dbConn *pgx.Conn, userID int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "get_all_snippet_by_userid", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet DBsnippet
		var code []byte
		err := rows.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description, &snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated)
		if err != nil {
			return nil, err
		}

		decompressedData, err := cmp.DecompressZSTD(code)
		if err != nil {
			return nil, err
		}

		snippet.Code = string(decompressedData)
		data = append(data, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil

}
