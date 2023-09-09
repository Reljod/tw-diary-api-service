package database

import (
	"context"
	"fmt"
	"os"

	"github.com/Reljod/tw-diary-api-service/config"
	"github.com/jackc/pgx/v5"
)

func BuildUrl(config config.ConfigSchema) string {
	db := config.Database
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		db.Engine, db.User, db.Password, db.Host, db.Port, db.Db)
}

func NewDatabase(config config.ConfigSchema) *pgx.Conn {
	url := BuildUrl(config)
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

var Conn *pgx.Conn = NewDatabase(config.Config)
