package repositories

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	DB = db
	log.Println("Database connected")
	return nil
}
