package main

import (
	"log"
	"net/http"
	"os"

	mux "github.com/gorilla/mux"
	controller "github.com/route-finder/ioan/routes/controller"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", controller.PrintHello).Methods("GET")
	r.HandleFunc("/routes", controller.GetRoutes).Methods("GET")

	log.Println("Serving on Port ..." + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
