package config

import (
	"database/sql"
	"fmt"
	"user-auth/internal/config/viper"

	_ "github.com/go-sql-driver/mysql"
)

// Creating a db variable of type *sql.DB
var db *sql.DB

func Connect() {
	// Getting the database credentials from config file
	username := viper.Get("database.username")
	password := viper.Get("database.password")
	port := viper.Get("database.port")
	database := viper.Get("database.database")

	localdb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", username, password, port, database))
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
