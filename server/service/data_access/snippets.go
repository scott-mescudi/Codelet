package dataaccess

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	cmp "github.com/scott-mescudi/codelet/shared/compression"
)

func AddSnippet(dbConn *pgxpool.Pool, userID int, language, description, title string, code string, private, favorite bool, tags []string, created time.Time, updated time.Time) error {
	compressed, err := cmp.CompressZSTD([]byte(code))
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(context.Background(), "INSERT INTO snippets(userid, language, title, code, description, private, tags, created, updated, favorite) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", userID, language, title, compressed, description, private, tags, created, updated, favorite)
	if err != nil {
		return err
	}

	return nil
}

func GetSnippetsByUserID(dbConn *pgxpool.Pool, userID, limit, offset int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE userid=$1 LIMIT $2 OFFSET $3", userID, limit, offset)
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

func GetAllSnippetsByUserID(dbConn *pgxpool.Pool, userID int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE userid=$1", userID)
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

func GetPublicSnippets(dbConn *pgxpool.Pool, limit, offset int) ([]DBsnippet, error) {
	var data []DBsnippet
	rows, err := dbConn.Query(context.Background(), "SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE private=false LIMIT $1 OFFSET $2", limit, offset)
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

func DeleteSnippet(dbConn *pgxpool.Pool, snippetID int) error {
	tx, err := dbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = dbConn.Exec(context.Background(), "DELETE FROM snippets WHERE id=$1", snippetID)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func GetSmallUserSnippets(dbConn *pgxpool.Pool, userID int) ([]SmallDBsnippet, error) {
	var data []SmallDBsnippet
	row, err := dbConn.Query(context.Background(), "SELECT id, language, title, favorite FROM snippets where userid=$1", userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var snippet SmallDBsnippet
		err = row.Scan(&snippet.ID, &snippet.Language, &snippet.Title, &snippet.Favorite)
		if err != nil {
			return nil, err
		}

		data = append(data, snippet)
	}

	return data, nil
}

func GetSnippetByIDAndUserID(dbConn *pgxpool.Pool, userID, snippetID int) (*DBsnippet, error) {
	var snippet DBsnippet
	var code []byte
	err := dbConn.QueryRow(context.Background(), "SELECT id, language, title, code, description, private, tags, created, updated, favorite FROM snippets WHERE userid=$1 AND id=$2", userID, snippetID).Scan(
		&snippet.ID, &snippet.Language, &snippet.Title, &code, &snippet.Description,
		&snippet.Private, &snippet.Tags, &snippet.Created, &snippet.Updated, &snippet.Favorite,
	)
	if err != nil {
		return nil, err
	}

	decompressedData, err := cmp.DecompressZSTD(code)
	if err != nil {
		return nil, err
	}

	snippet.Code = string(decompressedData)
	return &snippet, nil
}

func UpdateUserSnippetByID(dbConn *pgxpool.Pool, snippetID int, language *string, title *string, code *string, favorite *bool, private *bool, tags *[]string, description *string) error {
	builder := strings.Builder{}
	args := 1
	builder.WriteString("UPDATE snippets SET")

	if language != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("language=$%d", args))
		args++
	}

	if title != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("title=$%d", args))
		args++
	}

	if code != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("code=$%d", args))
		args++
	}

	if favorite != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("favorite=$%d", args))
		args++
	}

	if private != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("private=$%d", args))
		args++
	}

	if tags != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("tags=$%d", args))
		args++
	}

	if description != nil {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("description=$%d", args))
		args++
	}

	if args == 1 {
		return errors.New("no fields to update")
	}

	builder.WriteString(" ")
	builder.WriteString(fmt.Sprintf("WHERE id=$%d", args))
	args++

	fmt.Println(builder.String())

	return nil

}

// UPDATE snippets SET language=$1 title=$2 code=$3 favorite=$4 private=$5 tags=$6 description=$7 WHERE id=$8
// UPDATE snippets SET language=$1 title=$2 code=$3 favorite=$4 private=$5 description=$6 WHERE id=$7
// UPDATE snippets SET language=$1 title=$2 code=$3 favorite=$4 private=$5 description=$6 WHERE id=$7
// UPDATE snippets SET language=$1 title=$2 code=$3 favorite=$4 private=$5 description=$6 WHERE id=$7
// UPDATE snippets SET language=$1 title=$2 code=$3 favorite=$4 private=$5 description=$6 WHERE id=$7

