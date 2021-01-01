package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/lestrrat-go/backoff"
	views "github.com/route-finder/ioan/routes/views"
)

// DoHTTPRequest perform a GET HTTP request on the given URL
func DoHTTPRequest(URL string, w http.ResponseWriter) ([]byte, int, error) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		WriteErrorResponse(resp.StatusCode, "Cannot send request to OSRM", w)
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("osrmURL: ", URL)
		log.Println("Error: ", err)

		WriteErrorResponse(resp.StatusCode, "Cannot read the response Body from OSRM", w)
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

// Setup the backoff policy
var policy = backoff.NewExponential(
	backoff.WithInterval(500*time.Millisecond), // base interval
	backoff.WithJitterFactor(0.05),             // 5% jitter
	backoff.WithMaxRetries(5),                  // If not specified, default number of retries is 10
)

// RetryHTTPRequest - a backoff algorithm to retry when the 500 Status Code was obtained
// This ensure that the request will be repeated when the endpoint
// reporst a Internal Server Error
// https://github.com/lestrrat-go/backoff
func RetryHTTPRequest(url string, w http.ResponseWriter) ([]byte, int, error) {
	background, cancel := policy.Start(context.Background())
	defer cancel()

	var statusCode = 500
	for backoff.Continue(background) {
		resp, statusCode, err := DoHTTPRequest(url, w)
		if err == nil {
			return resp, statusCode, nil
		}
	}

	return nil, statusCode, errors.New(`Tried very hard, but no luck`)
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
func WriteErrorResponse(statusCode int, message string, w http.ResponseWriter) {
	response := views.Response{}

	response.Code = statusCode
	response.Body = message

	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}
