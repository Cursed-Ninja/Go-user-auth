// All api handlers are defined here

package controllers

import (
	"log"
	"net/http"
	"user-auth/internal/config/viper"
	oAuth "user-auth/internal/googleOauth"
	"user-auth/internal/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Session related variables
var (
	store       = sessions.NewCookieStore([]byte(viper.Get("session.secret")))
	sessionName = viper.Get("session.name")
)

// Handler for login
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing form to get details of the user from login page

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := models.GetUser(email)

	if err != nil {
		log.Println(err)
		http.Error(w, "Username entered does not exist", http.StatusNotFound)
		return
	}

	// If user is not registered using email and password, then return error
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

	// Setting the session
	session, _ := store.Get(r, sessionName)
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

	// If email or password is empty, then return error
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

	// If user already exists, then return error
	if err != nil {
		log.Println(err)
		http.Error(w, "User already exists, try logging in", http.StatusBadRequest)
		return
	}

	// Setting the session
	session, _ := store.Get(r, sessionName)

	session.Values["authenticated"] = true
	session.Values["email"] = newUser.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/edit-details", http.StatusSeeOther)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
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

	// Creating a new user object to update the user details
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
	session.Values["email"] = newUser.Email
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Handler for Google Oauth
func GoogleOauthHandler(w http.ResponseWriter, r *http.Request) {
	url := oAuth.GoogleOauth()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handler for Google Oauth callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oAuth.Callback(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if user already exists
	user, err := models.GetUser(userInfo.Email)
	user.Password = nil
	isNewUser := false

	if err != nil {
		// If user does not exist, then create a new user
		log.Println(err)
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

	// If user is not registered using google, then return error
	if !isNewUser && user.SignInMethod != "google" {
		log.Println(err)
		http.Redirect(w, r, "/login?error=409", http.StatusSeeOther)
		return
	}
	
	// Setting the session
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	// Redirecting to the appropriate page
	var url string
	if isNewUser {
		url = "/edit-details"
	} else {
		url = "/profile"
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Emptying the session
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
