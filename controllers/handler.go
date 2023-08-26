package controllers

import (
	"log"
	"net/http"
	"user-auth/config/viper"
	oAuth "user-auth/googleOauth"
	"user-auth/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	store       = sessions.NewCookieStore([]byte(viper.Get("session.secret")))
	sessionName = viper.Get("session.name")
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	err := r.ParseMultipartForm(32 << 20)
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
		http.Error(w, "Username entered does not exist", http.StatusNotFound)
		return
	}

	if user.SignInMethod != "email" {
		log.Println(err)
		http.Error(w, "Incorrect method. Try logging in using Google", http.StatusConflict)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		log.Println(err)
		http.Error(w, "Password is incorrect", http.StatusBadRequest)
		return
	}

	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}

	newUser.Name = r.FormValue("name")
	newUser.Email = r.FormValue("email")
	newUser.Phone = r.FormValue("phone")
	newUser.SignInMethod = "email"

	password := r.FormValue("password")

	if newUser.Email == "" || password == "" {
		log.Println("Please provide username and password")
		http.Error(w, "Please provide username and password", http.StatusBadRequest)
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
		http.Error(w, "User already exists, try logging in", http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, sessionName)

	session.Values["authenticated"] = true
	session.Values["email"] = newUser.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/edit-details", http.StatusSeeOther)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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
	previousEmail := session.Values["email"].(string)
	if session.Values["googleOauth"].(bool) {
		newUser.Email = previousEmail
	}
	err = newUser.UpdateUser(previousEmail)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func GoogleOauthHandler(w http.ResponseWriter, r *http.Request) {
	url := oAuth.GoogleOauth()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oAuth.Callback(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := models.GetUserDetails(userInfo.Email)
	isNewUser := false
	if err != nil {
		newUser := models.User{}
		newUser.Email = userInfo.Email
		newUser.Name = userInfo.Name
		newUser.SignInMethod = "google"
		err = newUser.RegisterUser()

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		isNewUser = true
	}

	if !isNewUser && user.SignInMethod != "google" {
		log.Println(err)
		http.Redirect(w, r, "/login?error=409", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	var url string
	if isNewUser {
		url = "/edit-details"
	} else {
		url = "/profile"
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
