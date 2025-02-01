package snippets

import (
	"github.com/jackc/pgx/v5"
)

type SnippetService struct {
	Db *pgx.Conn
}

type Snippet struct {
	Language    string   `json:"language"`
	Title       string   `json:"title"`
	Code        string   `json:"code"`
	Favorite    bool     `json:"favorite"`
	Private     bool     `json:"private"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}
