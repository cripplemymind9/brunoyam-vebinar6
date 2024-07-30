package storage

import (
	"fmt"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(conn  *pgx.Conn) (*PostgresStorage, error) {
	return &PostgresStorage{
		conn: conn,
	}, nil
}

func Migrations(dbAddr, migrationsPath string) error {
	migratePath := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(migratePath, dbAddr)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}