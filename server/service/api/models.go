package api

import "github.com/jackc/pgx/v5"

type Server struct {
	Db *pgx.Conn
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignup struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type AddSnippet struct {
	Userid      int      `json:"userid"`
	Language    string   `json:"language"`
	Title       string   `json:"title"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type DeleteSnippet struct {
	Id     int `json:"id"`
	Userid int `json:"userid"`
}
