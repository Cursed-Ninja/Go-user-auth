// All static file handlers are defined here

package controllers

import (
	"log"
	"net/http"
	"user-auth/internal/models"
	"user-auth/internal/templates"
)

// helper function for rendering templates
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
	// Check if user is authenticated
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := session.Values["email"].(string)
	isGoogleOauth := session.Values["googleOauth"].(bool)
	user, err := models.GetUser(email)
	user.Password = nil

	if err != nil {
		log.Println(err)
		http.Error(w, "User does not exists", http.StatusBadRequest)
		return
	}

	data := struct {
		models.User
		GoogleOauth bool
	}{
		User:        user,
		GoogleOauth: isGoogleOauth,
	}

	// Prevents caching of the page
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	renderTemplate(w, "edit-details", data)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := session.Values["email"].(string)

	user, err := models.GetUser(email)
	user.Password = nil

	if err != nil {
		log.Println(err, email, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Prevents caching of the page
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	renderTemplate(w, "profile", user)
}
