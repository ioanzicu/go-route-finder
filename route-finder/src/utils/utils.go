package utils

import (
	"encoding/json"
	views "github.com/route-finder/ioan/routes/views"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

func DoHTTPRequest(URL string, w http.ResponseWriter) ([]byte, error) {
	response := views.Response{}

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
