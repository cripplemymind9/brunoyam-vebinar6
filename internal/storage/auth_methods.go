package storage

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"context"
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
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
	var hashedPassword string
	for rows.Next() {
		err := rows.Scan(
			&uid,
			&hashedPassword,
		)

		if err != nil {
			return -1, fmt.Errorf("invalid login credentials")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		return -1, fmt.Errorf("invalid login credentials")
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