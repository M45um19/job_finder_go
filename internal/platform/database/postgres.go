package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(DBUrl string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), DBUrl)

	if err != nil {
		log.Fatal("Database connection fail", err)
	}

	return pool
}
