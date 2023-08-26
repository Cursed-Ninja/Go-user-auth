package routes

import (
	"net/http"
	"user-auth/controllers"

	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/register", controllers.RegisterUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/register", controllers.RegisterUser).Methods(http.MethodGet)
	router.HandleFunc("/login", controllers.LoginUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", controllers.LoginUser).Methods(http.MethodGet)
	router.HandleFunc("/update-profile", controllers.UpdateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/edit-details", controllers.EditUser).Methods(http.MethodGet)
	router.HandleFunc("/profile", controllers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/oauth/google/login", controllers.GoogleOauthLoginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauth/callback/login", controllers.CallbackLoginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauth/google/register", controllers.GoogleOauthRegisterHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauth/callback/register", controllers.CallbackRegisterHandler).Methods(http.MethodGet)
	router.HandleFunc("/logout", controllers.LogoutHandler).Methods(http.MethodGet)
	router.HandleFunc("/", controllers.LoginUser).Methods(http.MethodGet)
}
