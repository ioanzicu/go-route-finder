package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	mux "github.com/gorilla/mux"
)

// Response struct
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
}

func printHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Code: http.StatusOK,
		Body: "Hello Horld!",
	}

	log.Printf("Hello, World!\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func getRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Code: http.StatusOK,
		Body: "Soon, here will be (an) EPIC ROUTE/S!",
	}

	log.Printf("EPIC ROUTES!\n")

	// We save to parse params:
	// src - one  				example: src=13.388860,52.517037
	// dst - one or many		example: dst=13.397634,52.529407&dst=13.428555,52.523219

	// Query()["tags"] will return a []string
	queryParameters := r.URL.Query()

	srcParams, srcOK := queryParameters["src"]
	dstParams, dstOK := queryParameters["dst"]

	// Accept just the input that consist of 'src' and 'dst' params
	if (!srcOK || len(srcParams[0]) < 1) || (!dstOK || len(dstParams[0]) < 1) {
		message := "Missing required query parameters: src and dst"

		log.Println(message)
		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check just for ONE 'src' query param
	if len(srcParams) > 1 {
		message := "Just one src param is allowed"

		log.Println(message)
		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	srcParams = strings.Split(srcParams[0], ",")
	srcParamsString := strings.Join(srcParams, ", ")
	response.Body = "Request Params: " + srcParamsString

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", printHello)
	r.HandleFunc("/routes", getRoutes)

	log.Println("Serving on Port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
