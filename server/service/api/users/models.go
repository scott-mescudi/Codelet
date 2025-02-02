package users

import (

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type UserService struct {
	Db     *pgxpool.Pool
	Logger zerolog.Logger
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
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
