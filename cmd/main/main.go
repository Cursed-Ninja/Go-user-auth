// Entry point

package main

import (
	"log"
	"net/http"
	"os"
	"user-auth/internal/config"
	"user-auth/internal/config/viper"
	"user-auth/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Setting up log file to log all the outputs

	if logToFile := viper.Get("log.log_to_file"); logToFile == "true" {
		logFile, err := os.Create("logfile.txt")
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	defer config.Disconnect()

	// Setting up mux router for routing
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)
	http.Handle("/", r)

	// Listning to Port 8080
	log.Println("Listening on port 8080")
	log.Println(http.ListenAndServe(":8080", nil))
}
