package storage

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
	"context"
	"fmt"
	"time"
)

func (ps *PostgresStorage) Login(input models.LoginUser) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT uid, password FROM users WHERE login = $1`
	if _, err := tx.Prepare(ctx, "select-uid-password", query); err != nil {
		return -1, err
	}

	row := tx.QueryRow(ctx, "select-uid-password", input.Login)

	var uid int
	var hashedPassword string
	if err := row.Scan(&uid, &hashedPassword); err != nil {
		return -1, fmt.Errorf("invalid login credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		return -1, fmt.Errorf("invalid login credentials")
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, err
	}

	return uid, nil
}

func (ps *PostgresStorage) Profile(claims models.Claims) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT * FROM users WHERE login = $1`
	if _, err := tx.Prepare(ctx, "select-user-by-login", query); err != nil {
		return models.User{}, err
	}

	rows, err := tx.Query(ctx, "select-user-by-login", claims.Login);
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(
			&user.UID,
			&user.Name,
			&user.Login,
			&user.Password,
		); err != nil {
			return models.User{}, err
		}
	}

	if err := rows.Err(); err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.User{}, err
	}

	return user, nil
}