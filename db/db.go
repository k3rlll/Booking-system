package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionDB() (*pgxpool.Pool, error) {

	connStr := "postgres://postgres:252566@localhost:5432/myapp"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	return pool, nil
}
