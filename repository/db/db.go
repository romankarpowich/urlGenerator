package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	Store *dbStore
)

type dbStore struct {
	*sql.DB
}

func InitDb() {
	Store = new(dbStore)
	connectionString := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	s, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	Store.DB = s
}
