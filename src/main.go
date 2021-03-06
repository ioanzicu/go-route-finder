package main

import (
	"log"
	"net/http"
	"time"

	mux "github.com/gorilla/mux"
	controller "github.com/route-finder/ioan/routes/controller"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", controller.PrintHello).Methods("GET", "OPTIONS")
	r.HandleFunc("/routes", controller.GetRoutes).Methods("GET", "OPTIONS")

	server := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Serving on Port 8000 ...")
	log.Fatal(server.ListenAndServe())
}
