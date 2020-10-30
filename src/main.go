package main

import (
	"log"
	"net/http"

	mux "github.com/gorilla/mux"
	controller "github.com/route-finder/ioan/routes/controller"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", controller.PrintHello).Methods("GET", "OPTIONS")
	r.HandleFunc("/routes", controller.GetRoutes).Methods("GET", "OPTIONS")

	log.Println("Serving on Port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
