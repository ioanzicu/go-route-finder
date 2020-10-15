package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// OSRMRoute struct
type OSRMRoute struct {
	Legs       []map[string]interface{} `json:"legs"`
	WeightName string                   `json:"weight_name"`
	Distance   float64                  `json:"distance"`
	Duration   float64                  `json:"duration"`
}

// OSRMResponse sturct
type OSRMResponse struct {
	Code      string                   `json:"code"`
	Waypoints []map[string]interface{} `json:"waypoints"`
	Routes    []OSRMRoute              `json:"routes"`
}

// ---------------

// Route struct
type Route struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

// OutputResponse struct
type OutputResponse struct {
	Source string  `json:"source"`
	Routes []Route `json:"routes"`
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

func doHTTPRequest(URL string, w http.ResponseWriter) ([]byte, error) {
	response := Response{}

	resp, err := http.Get(URL)
	if err != nil {
		response.Body = "Cannot send request to OSRM"
		response.Code = http.StatusBadRequest
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Body = "Cannot read the response Body from OSRM"
		response.Code = http.StatusBadRequest
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return nil, err
	}

	return body, nil
}

func getRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	const responseCodeOK = "Ok"

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

	log.Printf("%v", routeRequest)

	// Validated src and dst
	sourceString := fmt.Sprintf("%g,%g", source.Latitude, source.Longitude)
	log.Println("source:", sourceString)

	destinationString := strings.Join(dstParams, ";")
	log.Println("destination:", destinationString)

	routes := []Route{}

	// Send request
	for _, dst := range dstParams {
		osrmURL := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s;%s?overview=false", sourceString, dst)

		reqBody, err := doHTTPRequest(osrmURL, w)
		if err != nil {
			return
		}

		osrmResp := OSRMResponse{}
		err = json.Unmarshal(reqBody, &osrmResp)

		if err != nil || osrmResp.Code != responseCodeOK {
			response.Body = "Cannot UNMARSHAL the response Body from OSRM or Code Response is not Ok"
			response.Code = http.StatusBadRequest
			log.Println("osrmURL: ", osrmURL)
			log.Println("Error: ", err)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Printf("%v", osrmResp.Routes[0])
		routes = append(routes, Route{
			Destination: dst,
			Duration:    osrmResp.Routes[0].Duration,
			Distance:    osrmResp.Routes[0].Distance,
		})
	}

	outResp := OutputResponse{
		Source: sourceString,
		Routes: routes,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(outResp)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", printHello)
	r.HandleFunc("/routes", getRoutes)

	log.Println("Serving on Port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
