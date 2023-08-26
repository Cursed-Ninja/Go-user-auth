// Entry point

package main

import (
	"log"
	"net/http"
	"os"
	"user-auth/internal/config"
	"user-auth/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Setting up log file to log all the outputs
	logFile, err := os.Create("logfile.txt")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	defer config.Disconnect()
	log.SetOutput(logFile)

	// Setting up mux router for routing
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)
	http.Handle("/", r)

	// Listning to Port 8080
	log.Println(http.ListenAndServe(":8080", nil))
}
