package dataaccess

import "time"

type DBsnippet struct {
	ID          int       `json:"id"`
	Language    string    `json:"language"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Private     bool      `json:"private"`
	Favorite    bool      `json:"favorite"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type SmallDBsnippet struct {
	ID       int    `json:"id"`
	Language string `json:"language"`
	Title    string `json:"title"`
	Favorite bool   `json:"favorite"`
}
