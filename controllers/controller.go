package controllers

import (
	"log"
	"net/http"
	"user-auth/config/viper"
	oAuth "user-auth/googleOauth"
	"user-auth/models"
	"user-auth/templates"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	user        models.User
	store       = sessions.NewCookieStore([]byte(viper.Get("session.secret")))
	sessionName = viper.Get("session.name")
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
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

	if user.SignInMethod != "email" {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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
	newUser.SignInMethod = "email"

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

	_, err = models.GetUserDetails(newUser.Email)

	if err == nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	newUser.Password = hashedPassword

	err = newUser.RegisterUser()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := session.Values["email"].(string)

	user, err := models.GetUserDetails(email)

	if err != nil {
		log.Println(err, email, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	renderTemplate(w, "profile", user)
}

func GoogleOauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oAuth.GoogleOauth("login")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackLoginHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oAuth.Callback(r, "login")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := models.GetUserDetails(userInfo.Email)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	if user.SignInMethod != "google" {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func GoogleOauthRegisterHandler(w http.ResponseWriter, r *http.Request) {
	url := oAuth.GoogleOauth("register")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackRegisterHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := oAuth.Callback(r, "register")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := models.GetUserDetails(userInfo.Email)

	if err == nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	user.Email = userInfo.Email
	user.Name = userInfo.Name
	user.SignInMethod = "google"

	err = user.RegisterUser()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, sessionName)

	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["googleOauth"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/edit-details", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.Templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := session.Values["email"].(string)
	isGoogleOauth := session.Values["googleOauth"].(bool)
	user, err := models.GetUserDetails(email)

	if err != nil {
		log.Println(err, email, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := struct {
		models.User
		GoogleOauth bool
	}{
		User:        user,
		GoogleOauth: isGoogleOauth,
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	renderTemplate(w, "edit-details", data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
