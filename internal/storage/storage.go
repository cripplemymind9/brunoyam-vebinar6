package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(conn  *pgx.Conn) (*PostgresStorage, error) {
	if err := CreatePostgresDB(conn); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		conn: conn,
	}, nil
}

func CreatePostgresDB(conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    	uid SERIAL PRIMARY KEY,
    	name TEXT,
    	login TEXT,
    	password TEXT
	);
	
	CREATE TABLE IF NOT EXISTS books (
		b_id SERIAL PRIMARY KEY,
		author TEXT,
		title TEXT,
		uid INTEGER
	);`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}