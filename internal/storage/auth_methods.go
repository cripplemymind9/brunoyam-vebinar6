package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
)

func (s *PostgresStorage) Login(input models.LoginUser) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT uid, password FROM users WHERE login = $1`
	rows, err := s.conn.Query(ctx, query, input.Login)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var uid int
	var password string
	for rows.Next() {
		err := rows.Scan(
			&uid,
			&password,
		)

		if err != nil {
			return -1, err
		}
	}

	if password != input.Password {
		return -1, fmt.Errorf("invalid data")
	}

	return uid, nil
}

func (s *PostgresStorage) Profile(claims models.Claims) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE login = $1`
	rows, err := s.conn.Query(ctx, query, claims.Login)
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