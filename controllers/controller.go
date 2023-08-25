package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"user-auth/models"

	"golang.org/x/crypto/bcrypt"
)

var user models.User

func LoginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := models.LoginUser(email)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(user.Password)
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}

	newUser.Name = r.FormValue("name")
	newUser.Email = r.FormValue("email")
	newUser.Phone = r.FormValue("phone")

	password := r.FormValue("password")

	if newUser.Email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser.Password = hashedPassword

	err = newUser.RegisterUser()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}

	newUser.Name = r.FormValue("name")
	newUser.Email = r.FormValue("email")
	newUser.Phone = r.FormValue("phone")

	err = newUser.UpdateUser()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	user, err := models.GetUserDetails(email)

	if err != nil {
		log.Println(err, email, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
