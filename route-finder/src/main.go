package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	mux "github.com/gorilla/mux"
)

// Response struct
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
}

// Source struct - Home Address
type Source struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Destination struct - Destination Address
type Destination struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RouteRequest struct for user input
type RouteRequest struct {
	Source      Source        `json:"src"`
	Destination []Destination `json:"dst"`
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

	// Get each parameters value as a []string
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

	log.Println(srcParams[0], (reflect.TypeOf(srcParams[0])))

	srcParams = strings.Split(srcParams[0], ",")

	// Make sure that 'src' has latitude and longitude provided
	if len(srcParams) != 2 {
		message := "Expect 'src' to have lattitude and longitude"

		log.Println(message)
		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Source
	srcParamsString := strings.Join(srcParams, ", ")
	response.Body = "Request Params: " + srcParamsString

	log.Println("latitude", srcParams[0], "longitude", srcParams[1],
		(reflect.TypeOf(srcParams[0])))

	tempSrcLatitude := srcParams[0]
	tempSrcLongitude := srcParams[1]

	log.Println("LOG", tempSrcLatitude, tempSrcLongitude)

	srcLatitude := 0.0

	if lat, err := strconv.ParseFloat(tempSrcLatitude, 64); err == nil {
		srcLatitude = lat
		log.Println("Lat: ", lat)
	} else {
		response.Code = http.StatusBadRequest
		response.Body = "Malformated param type (float64)"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	srcLongitude := 0.0

	if lon, err := strconv.ParseFloat(tempSrcLongitude, 64); err == nil {
		srcLongitude = lon

		log.Println("Long: ", lon)
	} else {
		response.Code = http.StatusBadRequest
		response.Body = "Malformated param type (float64)"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	source := Source{
		Latitude:  srcLatitude,
		Longitude: srcLongitude,
	}

	// Destination

	dstLatitude := 0.0
	dstLongitude := 0.0

	var destinationSlice []Destination

	for index, destination := range dstParams {
		log.Println("Destination: ", index, " ", destination, len(destination))

		arrDestination := strings.Split(destination, ",")

		// Make sure that 'dst' has latitude and longitude provided
		if len(arrDestination) != 2 {
			message := "Expect 'dst' to have lattitude and longitude"

			log.Println(message)
			response.Code = http.StatusBadRequest
			response.Body = message

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Println("Request Params: "+destination, "[]string", arrDestination)
		log.Println("latitude", arrDestination[0], "longitude", arrDestination[1],
			(reflect.TypeOf(arrDestination[0])))

		tempSrcLatitude := arrDestination[0]
		tempSrcLongitude := arrDestination[1]

		if lat, err := strconv.ParseFloat(tempSrcLatitude, 64); err == nil {
			dstLatitude = lat
			log.Println("Lat: ", lat)
		} else {
			response.Body = "Malformated param type (float64)"
			response.Code = http.StatusBadRequest

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if long, err := strconv.ParseFloat(tempSrcLongitude, 64); err == nil {
			dstLongitude = long

			log.Println("Long: ", long)
		} else {
			response.Body = "Malformated param type (float64)"
			response.Code = http.StatusBadRequest

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		tempDst := Destination{
			Latitude:  dstLatitude,
			Longitude: dstLongitude,
		}

		destinationSlice = append(destinationSlice, tempDst)
	}

	routeRequest := RouteRequest{
		Source:      source,
		Destination: destinationSlice,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routeRequest)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", printHello)
	r.HandleFunc("/routes", getRoutes)

	log.Println("Serving on Port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
