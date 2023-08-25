package main

import (
	"log"
	"net/http"
	"user-auth/config"
	"user-auth/routes"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	logFile, err := os.Create("logfile.txt")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)
	http.Handle("/", r)
	log.Println(http.ListenAndServe(":8080", nil))
	defer config.Disconnect()
}
