package dataaccess

import (
	"context"
	"github.com/jackc/pgx/v5"
	cmp "github.com/scott-mescudi/codelet/shared/compression"
	"time"
)

func AddSnippet(dbConn *pgx.Conn, userID int, language, description, title string, code string, private, favorite bool, tags []string, created time.Time, updated time.Time) error {
	compressed, err := cmp.CompressZSTD([]byte(code))
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(context.Background(), "add_snippet", userID, language, title, compressed, description, private, tags, created, updated, favorite)
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
		err := rows.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description, &snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated, &snippet.Favorite)
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
		err := rows.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description, &snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated, &snippet.Favorite)
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

func GetPublicSnippets(dbConn *pgx.Conn, limit, offset int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "get_public_snippet", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet DBsnippet
		var code []byte
		err := rows.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description, &snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated, &snippet.Favorite)
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

func DeleteSnippet(dbConn *pgx.Conn, snippetID int) error {
	tx, err := dbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = dbConn.Exec(context.Background(), "delete_snippet", snippetID)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}
