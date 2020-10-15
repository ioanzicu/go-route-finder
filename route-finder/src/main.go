package main

import (
	"log"
	"net/http"

	mux "github.com/gorilla/mux"
	controller "github.com/route-finder/ioan/routes/controller"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", controller.PrintHello).Methods("GET")
	r.HandleFunc("/routes", controller.GetRoutes).Methods("GET")

	log.Println("Serving on Port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
