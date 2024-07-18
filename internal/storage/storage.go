package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, connStr string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := CreatePostgresDB(ctx, conn); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		conn: conn,
	}, nil
}

func CreatePostgresDB(ctx context.Context, conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    	uid SERIAL PRIMARY KEY,
    	name TEXT,
    	login TEXT,
    	password TEXT
	);`

	_, err := conn.Exec(ctx, query)
	return err
}