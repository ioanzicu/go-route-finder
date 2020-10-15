package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	utils "github.com/route-finder/ioan/routes/utils"
	views "github.com/route-finder/ioan/routes/views"
)

// PrintHello send an EPIC 'Hello World!' response JSON message
func PrintHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := views.Response{
		Code: http.StatusOK,
		Body: "Hello World!",
	}

	log.Printf("Hello, World!\n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRoutes parse the GET parameters: 'src' & 'dst'
// Each parameter contains latitude and longitude coordinates
// Then, it do a request to the Open Street Maps API
// Result is reshaped, sorted and send to the user in JSON format
func GetRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	const responseCodeOK = "Ok"

	response := views.Response{
		Code: http.StatusOK,
		Body: "Soon, here will be (an) EPIC ROUTE/S!",
	}

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

		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check just for ONE 'src' query param
	if len(srcParams) > 1 {
		message := "Just one src param is allowed"

		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	srcParams = strings.Split(srcParams[0], ",")

	// Make sure that 'src' has latitude and longitude provided
	if len(srcParams) != 2 {
		message := "Expect 'src' to have lattitude and longitude"

		response.Code = http.StatusBadRequest
		response.Body = message

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Source
	srcParamsString := strings.Join(srcParams, ", ")
	response.Body = "Request Params: " + srcParamsString

	tempSrcLatitude := srcParams[0]
	tempSrcLongitude := srcParams[1]

	srcLatitude := 0.0

	if lat, err := strconv.ParseFloat(tempSrcLatitude, 64); err == nil {
		srcLatitude = lat
	} else {
		response.Code = http.StatusBadRequest
		response.Body = "Malformated param type (float64)"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	srcLongitude := 0.0

	if long, err := strconv.ParseFloat(tempSrcLongitude, 64); err == nil {
		srcLongitude = long
	} else {
		response.Code = http.StatusBadRequest
		response.Body = "Malformated param type (float64)"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	source := views.Source{
		Latitude:  srcLatitude,
		Longitude: srcLongitude,
	}

	// Destination
	dstLatitude := 0.0
	dstLongitude := 0.0

	var destinationSlice []views.Destination

	for _, destination := range dstParams {
		// log.Println("Destination: ", index, " ", destination, len(destination))

		arrDestination := strings.Split(destination, ",")

		// Make sure that 'dst' has latitude and longitude provided
		if len(arrDestination) != 2 {
			message := "Expect 'dst' to have lattitude and longitude"

			response.Code = http.StatusBadRequest
			response.Body = message

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		tempSrcLatitude := arrDestination[0]
		tempSrcLongitude := arrDestination[1]

		if lat, err := strconv.ParseFloat(tempSrcLatitude, 64); err == nil {
			dstLatitude = lat
		} else {
			response.Body = "Malformated param type (float64)"
			response.Code = http.StatusBadRequest

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if long, err := strconv.ParseFloat(tempSrcLongitude, 64); err == nil {
			dstLongitude = long
		} else {
			response.Body = "Malformated param type (float64)"
			response.Code = http.StatusBadRequest

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		tempDst := views.Destination{
			Latitude:  dstLatitude,
			Longitude: dstLongitude,
		}

		destinationSlice = append(destinationSlice, tempDst)
	}

	routeRequest := views.RouteRequest{
		Source:      source,
		Destination: destinationSlice,
	}

	log.Printf("%v", routeRequest)

	// Validated 'src' and 'dst'
	sourceString := fmt.Sprintf("%g,%g", source.Latitude, source.Longitude)

	routes := []views.Route{}

	// Send request
	for _, dst := range dstParams {
		osrmURL := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s;%s?overview=false", sourceString, dst)

		reqBody, err := utils.DoHTTPRequest(osrmURL, w)
		if err != nil {
			return
		}

		osrmResp := views.OSRMResponse{}
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

		routes = append(routes, views.Route{
			Destination: dst,
			Duration:    osrmResp.Routes[0].Duration,
			Distance:    osrmResp.Routes[0].Distance,
		})
	}

	// Sort the routes
	routes = utils.SortSlice(routes)

	outResp := views.OutputResponse{
		Source: sourceString,
		Routes: routes,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(outResp)
}
