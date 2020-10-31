package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	views "github.com/route-finder/ioan/routes/views"
)

// DoHTTPRequest perform a GET HTTP request on the given URL
func DoHTTPRequest(URL string, w http.ResponseWriter) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		WriteErrorResponse("Cannot send request to OSRM", w)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		WriteErrorResponse("Cannot read the response Body from OSRM", w)
		return nil, err
	}

	return body, nil
}

// SortSlice sort the parameter slice by Duration field
// If There are some Duration fields equal, then, the
// slice will be also sorted by the Distance field
func SortSlice(slice []views.Route) []views.Route {
	equalDistance := false

	// Sort by driving time
	sort.Slice(slice, func(i, j int) bool {
		if slice[i].Duration == slice[j].Duration {
			equalDistance = true
		}
		return slice[i].Duration < slice[j].Duration
	})

	// Sort by distance (if driving time is equal)
	if equalDistance {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Distance < slice[j].Distance
		})
	}

	return slice
}

// ParseToFloat64 parse the string parameter to float64
func ParseToFloat64(strNumber string, w http.ResponseWriter) (float64, error) {
	srcLatitude := 0.0
	response := views.Response{}

	if lat, err := strconv.ParseFloat(strNumber, 64); err == nil {
		srcLatitude = lat
	} else {
		response.Code = http.StatusBadRequest
		response.Body = "Malformated param type (float64)"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return 0.0, err
	}

	return srcLatitude, nil
}

// WriteErrorResponse send to the client a JSON object
// with the Status Code and Error Message
func WriteErrorResponse(message string, w http.ResponseWriter) {
	response := views.Response{}

	response.Code = http.StatusBadRequest
	response.Body = message

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}
