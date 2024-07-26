package main

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/server"
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/storage"
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/config"
	"github.com/jackc/pgx/v5"
	"context"
	"log"
)

func main() {
	config := config.ReadConfig()

	conn, err := initDB(config.DBAddr)
	if err != nil {
		log.Fatal(err)
	}


	// TODO - "postgres://postgres:Gogaminoga2019&@localhost:5432/postgres"
	store, err := storage.NewPostgresStorage(conn)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(":8080", store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initDB(addr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}