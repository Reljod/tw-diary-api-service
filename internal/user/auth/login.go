package auth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

type InvalidUsernameOrPasswordError struct{}

func (e *InvalidUsernameOrPasswordError) Error() string {
	return "Invalid username or password"
}

func (auth *SimpleSessionBasedAuth) Login(username string, password string) error {

	isCredsEmpty := strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""

	if err := auth.verifyCredentials(username, password); isCredsEmpty || err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return &InvalidUsernameOrPasswordError{}
	}

	return nil
}

func (auth *SimpleSessionBasedAuth) verifyCredentials(username string, password string) error {
	var user, hashedPassword string
	err := auth.Db.QueryRow(context.Background(),
		"SELECT username, password FROM public.accounts WHERE username = $1",
		username).Scan(&user, &hashedPassword)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return &AuthError{}
	}

	if auth.PasswordManager.IsMatch(password, hashedPassword) != nil || err == pgx.ErrNoRows {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return &InvalidUsernameOrPasswordError{}
	}

	return nil
}
