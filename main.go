package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mux "github.com/gorilla/mux"
	controller "github.com/route-finder/ioan/routes/controller"
)

func main() {

	r := mux.NewRouter()
	// Method OPTIONS allow to setup cors in the handler function
	r.HandleFunc("/", controller.PrintHello).Methods("GET", "OPTIONS")
	r.HandleFunc("/routes", controller.GetRoutes).Methods("GET", "OPTIONS")

	server := &http.Server{
		Handler: r,
		Addr:    GetPort(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Serving on Port " + GetPort())
	log.Fatal(server.ListenAndServe())
}

// GetPort get the Port from the environment
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
