package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
)

type UsernameAlreadyTaken struct{}

func (e *UsernameAlreadyTaken) Error() string {
	return "Username is already taken"
}

func (auth *SimpleSessionBasedAuth) Register(username string, password string) (string, error) {

	hashedPassword, err := auth.PasswordManager.Hash(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return "", err
	}

	tx, err := auth.Db.Begin(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return "", err
	}

	defer tx.Rollback(context.Background())

	var id string
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO public.accounts (username, password) values ($1, $2) RETURNING id",
		username, hashedPassword).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", &UsernameAlreadyTaken{}
		}

		return "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return "", err
	}

	return id, nil
}
