package snippets

import "github.com/jackc/pgx/v5"

type SnippetService struct {
	Db *pgx.Conn
}

type AddSnippet struct {
	Language    string   `json:"language"`
	Title       string   `json:"title"`
	Code        string   `json:"code"`
	Private     bool     `json:"private"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type DeleteSnippet struct {
	Id     int `json:"id"`
	Userid int `json:"userid"`
}
