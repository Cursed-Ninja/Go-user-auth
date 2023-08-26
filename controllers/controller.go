package controllers

import (
	"log"
	"net/http"
	"user-auth/models"
	"user-auth/templates"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	user  models.User
	store = sessions.NewCookieStore([]byte("your-secret-key"))
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
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

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
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

	session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = true
	session.Values["email"] = newUser.Email
	session.Save(r, w)

	http.Redirect(w, r, "/edit-details", http.StatusSeeOther)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

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

	err = newUser.UpdateUser(previousEmail)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

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

	renderTemplate(w, "profile", user)
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
	session, _ := store.Get(r, "session-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := session.Values["email"].(string)
	log.Println(email)
	user, err := models.GetUserDetails(email)

	if err != nil {
		log.Println(err, email, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	renderTemplate(w, "edit-details", user)
}
