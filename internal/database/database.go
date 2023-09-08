package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewDatabase() *pgx.Conn {
	url := "postgres://user:password@localhost:5433/tw-diary"
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

var Conn *pgx.Conn = NewDatabase()
