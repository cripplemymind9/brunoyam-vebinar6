package storage

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"

	"context"
	"time"

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

func (s *PostgresStorage) GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.conn.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UID,
			&user.Name,
			&user.Login,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *PostgresStorage) GetUser(uid int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE uid = $1`
	rows, err := s.conn.Query(ctx, query, uid)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		err := rows.Scan(
			&user.UID,
			&user.Name,
			&user.Login,
			&user.Password,
		)

		if err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (s *PostgresStorage) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `INSERT INTO users (uid, name, login, password) VALUES ($1, $2, $3, $4)`
	_, err := s.conn.Exec(ctx, query, user.UID, user.Name, user.Login, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) UpdateUser(uid int, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `UPDATE users SET name = $1, login = $2, password = $3 WHERE uid = $4`
	_, err := s.conn.Exec(ctx, query, user.Name, user.Login, user.Password, uid)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) DeleteUser(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM users WHERE uid = $1`
	_, err := s.conn.Exec(ctx, query, uid)
	if err != nil {
		return err
	}

	return nil
}
