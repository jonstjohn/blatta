package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPoolFromUrl(url string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), url)
}
