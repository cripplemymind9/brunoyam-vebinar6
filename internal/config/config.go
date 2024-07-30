package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr 		string
	DBAddr 		string
	MPath     	string
}

const (
	defaultAddr  = ":8080"
	defaultDbDSN = "postgres://postgres:Gogaminoga2019&@localhost:5432/postgres?sslmode=disable"
	defaultMigratePath = "migrations"
)

func ReadConfig() Config {
	var addr 		string
	var dbAddr 		string
	var migratePath string
	flag.StringVar(&addr, "addr", defaultAddr, "Server address") // mani.exe -help
	flag.StringVar(&dbAddr, "db", defaultDbDSN, "Database connection addres")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")
	flag.Parse()

	if temp := os.Getenv("SERVER_ADDR"); temp != "" {
		if addr == defaultAddr {
			addr = temp
		}
	}

	if temp := os.Getenv("DB_DSN"); temp != "" {
		if dbAddr == defaultDbDSN {
			dbAddr = temp
		}
	}

	if temp := os.Getenv("MIGRATE_PATH"); temp != "" {
		if migratePath == defaultMigratePath {
			migratePath = temp
		}
	}

	return Config {
		Addr:   	addr,
		DBAddr: 	dbAddr,
		MPath:     	migratePath,
	}
}
