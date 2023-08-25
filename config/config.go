package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() {
	localdb, err := sql.Open("mysql", "admin:@Admin#123@tcp(localhost:3306)/test")
	if err != nil {
		panic("Failed to connect to database")
	}

	db = localdb
}

func Disconnect() {
	db.Close()
}

func GetDB() *sql.DB {
	return db
}
