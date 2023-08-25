package routes

import (
	"net/http"
	"user-auth/controllers"

	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/register", controllers.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/login", controllers.LoginUser).Methods(http.MethodPost)
	router.HandleFunc("/update-profile", controllers.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/profile", controllers.GetUser).Methods(http.MethodGet)
}
