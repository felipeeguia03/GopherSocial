package main

import (
	"log"

	"github.com/felipeeguia03/vol7/internal/db"
	"github.com/felipeeguia03/vol7/internal/env"
	"github.com/felipeeguia03/vol7/internal/store"
)

func main() {

	dsn := env.GetString("DB_ADDR", "postgres://root:root@localhost/vol7?sslmode=disable")
	conn, err := db.New(dsn, 10, 10, "3m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)

	if err := db.Seed(store, conn); err != nil {
		log.Fatal(err)
	}

}
