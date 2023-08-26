package models

import (
	"database/sql"
	"user-auth/config"
)

var db *sql.DB

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

func (u *User) RegisterUser() error {
	query := "INSERT INTO users (name, email, password, phone, method) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, u.Name, u.Email, u.Password, u.Phone, u.SignInMethod)
	return err
}

func LoginUser(email string) (User, error) {
	var user User
	query := "SELECT name, email, password, phone, method FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.Password, &user.Phone, &user.SignInMethod)
	return user, err
}

func GetUserDetails(email string) (User, error) {
	var user User
	query := "SELECT name, email, phone, method FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.Phone, &user.SignInMethod)
	return user, err
}

func (u *User) UpdateUser(previousEmail string) error {
	query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE email = ?"
	_, err := db.Exec(query, u.Name, u.Email, u.Phone, previousEmail)
	return err
}
