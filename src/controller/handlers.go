package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	utils "github.com/route-finder/ioan/routes/utils"
	views "github.com/route-finder/ioan/routes/views"
)

// PrintHello send an EPIC 'Hello World!' response JSON message
func PrintHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	const responseCodeOK = "Ok"

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
		utils.WriteErrorResponse(http.StatusBadRequest, "Missing required query parameters: src and/or dst", w)
		return
	}

	// Check just for ONE 'src' query param
	if len(srcParams) > 1 {
		utils.WriteErrorResponse(http.StatusBadRequest, "Just one `src` param is allowed", w)
		return
	}

	srcParams = strings.Split(srcParams[0], ",")

	// Make sure that 'src' has latitude and longitude provided
	if len(srcParams) != 2 {
		utils.WriteErrorResponse(http.StatusBadRequest, "Expect `src` to have lattitude and longitude", w)
		return
	}

	// Source
	tempSrcLatitude := srcParams[0]
	tempSrcLongitude := srcParams[1]

	srcLatitude, err := utils.ParseToFloat64(tempSrcLatitude, w)
	if err != nil {
		return
	}

	srcLongitude, err := utils.ParseToFloat64(tempSrcLongitude, w)
	if err != nil {
		return
	}

	source := views.Source{
		Latitude:  srcLatitude,
		Longitude: srcLongitude,
	}

	// Destination
	var destinationSlice []views.Destination

	for _, destination := range dstParams {
		arrDestination := strings.Split(destination, ",")

		// Make sure that 'dst' has latitude and longitude provided
		if len(arrDestination) != 2 {
			utils.WriteErrorResponse(http.StatusBadRequest, "Expect 'dst' to have lattitude and longitude", w)
			return
		}

		tempDstLatitude := arrDestination[0]
		tempDstLongitude := arrDestination[1]

		dstLatitude, err := utils.ParseToFloat64(tempDstLatitude, w)
		if err != nil {
			return
		}

		dstLongitude, err := utils.ParseToFloat64(tempDstLongitude, w)
		if err != nil {
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

		reqBody, statusCode, err := utils.DoHTTPRequest(osrmURL, w)
		if err != nil {
			return
		}
		// Retry the request if the response code is 500
		if statusCode == http.StatusInternalServerError {
			reqBody, statusCode, err = utils.RetryHTTPRequest(osrmURL, w)
			if err != nil {
				return
			}
		}

		osrmResp := views.OSRMResponse{}
		err = json.Unmarshal(reqBody, &osrmResp)

		if err != nil || osrmResp.Code != responseCodeOK {
			utils.WriteErrorResponse(statusCode, "Cannot UNMARSHAL the response Body from OSRM or Code Response is not Ok", w)
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
