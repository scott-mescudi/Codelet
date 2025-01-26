package dataaccess

import "time"

type DBsnippet struct {
	ID          int       `json:"id"`
	Language    string    `json:"language"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Private     bool      `json:"private"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
