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
	CreateSession() (*Session, error)
	IsSessionValid(sessionId string) (bool, error)
}

type SimpleSessionBasedAuth struct {
	Db              *pgx.Conn
	PasswordManager PasswordManager
	SessionHandler  SessionHandler
}

func (auth *SimpleSessionBasedAuth) CreateSession() (*Session, error) {
	return auth.SessionHandler.CreateNew()
}

func (auth *SimpleSessionBasedAuth) IsSessionValid(sessionId string) (bool, error) {
	return auth.SessionHandler.IsValid(sessionId)
}
