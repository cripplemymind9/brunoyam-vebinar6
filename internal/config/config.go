package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr 	string
	DBAddr 	string
}

const (
	defaultAddr  = ":8080"
	defaultDbDSN = "postgres://postgres:Gogaminoga2019&@localhost:5432/postgres"
)

func ReadConfig() Config {
	var addr string
	var dbAddr string
	flag.StringVar(&addr, "addr", defaultAddr, "Server address") // mani.exe -help
	flag.StringVar(&dbAddr, "db", defaultDbDSN, "Database connection addres")
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

	return Config {
		Addr:   addr,
		DBAddr: dbAddr,
	}
}