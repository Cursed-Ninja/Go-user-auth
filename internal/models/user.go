// All the database related functions are defined here
// This file is responsible for creating, fetching and updating the user in the database

package models

import (
	"database/sql"
	"user-auth/internal/config"
)

// Create a global variable to hold the database connection pool
var db *sql.DB

// User struct defines the structure of the user in the database
type User struct {
	Name         string
	Email        string
	Password     []byte
	Phone        string
	SignInMethod string
}

func init() {
	config.Connect()
	db = config.GetDB()
}

// Creates a new user in database
func (u *User) RegisterUser() error {
	query := "INSERT INTO users (name, email, password, phone, method) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, u.Name, u.Email, u.Password, u.Phone, u.SignInMethod)
	return err
}

// Fetches the user from the database
func GetUser(email string) (User, error) {
	var user User
	query := "SELECT name, email, password, phone, method FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.SignInMethod)
	return user, err
}

// Updates the user in the database
func (u *User) UpdateUser(previousEmail string) error {
	query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE email = ?"
	_, err := db.Exec(query, u.Name, u.Email, u.Phone, previousEmail)
	return err
}
