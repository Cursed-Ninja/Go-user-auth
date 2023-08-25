package models

import (
	"database/sql"
	"user-auth/config"
)

var db *sql.DB

type User struct {
	Name     string
	Email    string
	Password []byte
	Phone    string
}

func init() {
	config.Connect()
	db = config.GetDB()
}

func CreateUserOAuth() {

}

func (u *User) RegisterUser() error {
	query := "INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, u.Name, u.Email, u.Password, u.Phone)
	return err
}

func LoginUser(email string) (User, error) {
	var user User
	query := "SELECT name, email, password, phone FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.Password, &user.Phone)
	return user, err
}

func GetUserDetails(email string) (User, error) {
	var user User
	query := "SELECT name, email, phone FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.Phone)
	return user, err
}

func (u *User) UpdateUser() error {
	query := "UPDATE users SET name=?, email=?, phone=? WHERE email=?"
	_, err := db.Exec(query, u.Name, u.Email, u.Phone, u.Email)
	return err
}
