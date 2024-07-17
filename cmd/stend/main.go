package main

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/server"
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/storage"

	"context"
	"log"
)

func main() {
	store, err := storage.NewPostgresStorage(context.TODO(), "postgres://postgres:Gogaminoga2019&@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(":8080", store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}