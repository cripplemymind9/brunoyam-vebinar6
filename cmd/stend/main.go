package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cripplemymind9/brunoyam-vebinar6/internal/config"
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/server"
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/storage"
	"github.com/jackc/pgx/v5"
)

func main() {
	config := config.ReadConfig()

	err := storage.Migrations(config.DBAddr, config.MPath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := initDB(config.DBAddr)
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewPostgresStorage(conn)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(":8080", store)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed{
			log.Fatal(err)
		}
	}()

	<- stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Server exiting")
}

func initDB(addr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}