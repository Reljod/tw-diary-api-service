package auth

import (
	"github.com/jackc/pgx/v5"
)

type AuthError struct{}

func (e *AuthError) Error() string {
	return "Authentication Error"
}

type Authenticator interface {
	Login(username string, password string) error
	Register(username string, password string) (string, error)
}

type SimpleSessionBasedAuth struct {
	Db              *pgx.Conn
	PasswordManager PasswordManager
}
