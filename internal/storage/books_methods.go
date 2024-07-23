package storage

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"context"
	"time"
)

func (ps *PostgresStorage) InsertBook(book models.Book, uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO books (author, title, uid) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(ctx, query, book.Author, book.Title, uid); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ps *PostgresStorage) InsertBooks(books []models.Book, uid int) error {
	for _, book := range books {
		if err := ps.InsertBook(book, uid); err != nil {
			return err
		}
	}
	
	return nil
}

func (ps *PostgresStorage) GetBooksByUserId(uid int) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := ps.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT * FROM books WHERE uid = $1`
	rows, err := tx.Query(ctx, query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.BookId,
			&book.Author,
			&book.Title,
			&book.UserID,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return books, nil
}