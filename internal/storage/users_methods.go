package storage

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"context"
	"time"
	"fmt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (ps *PostgresStorage) GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	tx, err := ps.conn.Begin(ctx)
	if err != nil{
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT * FROM users`
	if _, err := tx.Prepare(ctx, "select-users", query); err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "select-users")
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return users, nil
}

func (ps *PostgresStorage) GetUser(uid int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT * FROM users WHERE uid = $1`
	if _, err := tx.Prepare(ctx, "select-user", query); err != nil {
		return models.User{}, err
	}
	
	row, err := tx.Query(ctx, "select-user", uid)
	if err != nil {
		return models.User{}, err
	}
	defer row.Close()

	var user models.User
	for row.Next() {
		err := row.Scan(
			&user.UID,
			&user.Name,
			&user.Login,
			&user.Password,
		)
		if err != nil {
			return models.User{}, err 
		}
	}

	if err := row.Err(); err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ps *PostgresStorage) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO users (uid, name, login, password) VALUES ($1, $2, $3, $4)`
	if _, err = tx.Exec(ctx, query, user.UID, user.Name, user.Login, hashedPassword); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ps *PostgresStorage) UpdateUser(uid int, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE users SET name = $1, login = $2, password = $3 WHERE uid = $4`
	if _, err = tx.Exec(ctx, query, user.Name, user.Login, hashedPassword, uid); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ps *PostgresStorage) DeleteUser(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM users WHERE uid = $1`
	if _, err = tx.Exec(ctx, query, uid); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ps *PostgresStorage) GetUserId(claims models.Claims) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT uid FROM users WHERE login = $1`
	if _, err := tx.Prepare(ctx, "select-uid", query); err != nil {
		return -1, err
	}

	row := tx.QueryRow(ctx, "select-uid", claims.Login)

	var uid int
	if err = row.Scan(&uid); err != nil {
		if err == pgx.ErrNoRows {
			return -1, fmt.Errorf("user not found")
		}
		return -1, err
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, err
	}

	return uid, nil
}