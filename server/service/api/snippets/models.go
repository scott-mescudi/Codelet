package snippets

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type SnippetService struct {
	Db     *pgxpool.Pool
	Logger zerolog.Logger
}

type UpdateSnippet struct {
	Language    *string   `json:"language"`
	Title       *string   `json:"title"`
	Code        *string   `json:"code"`
	Favorite    *bool     `json:"favorite"`
	Private     *bool     `json:"private"`
	Tags        *[]string `json:"tags"`
	Description *string   `json:"description"`
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
