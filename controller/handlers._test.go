package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	views "github.com/route-finder/ioan/routes/views"
)

func TestHelloWorld(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	PrintHello(recorder, request)

	responseRecorder := recorder.Result()

	var respMessage views.Response
	err := json.NewDecoder(responseRecorder.Body).Decode(&respMessage)
	if err != nil {
		t.Errorf("%s", err)
	}

	expectedResponse := views.Response{
		Code: http.StatusOK,
		Body: "Hello World! GitHub + Jenkins",
	}

	if (respMessage.Body != expectedResponse.Body) && (expectedResponse.Body != "") {
		t.Errorf("`%v` failed, got %v, expected %v", "TestHelloWorld", responseRecorder.Body, expectedResponse.Body)
	}

	if responseRecorder.StatusCode != expectedResponse.Code {
		t.Errorf("`%v` failed, got %v, expected %v", "TestHelloWorld", responseRecorder.StatusCode, expectedResponse.Code)
	}
}

func TestGetRoutes(t *testing.T) {
	testCases := map[string]struct {
		params     map[string]string
		message    string
		statusCode int
	}{
		"good params": {
			map[string]string{
				"src": "13.38886,52.517037",
				"dst": "13.397634,52.529407",
			},
			"",
			http.StatusOK,
		},
		"without params": {
			map[string]string{},
			"Missing required query parameters: src and/or dst",
			http.StatusBadRequest,
		},
		"malformated `src` params, no latitude": {
			map[string]string{
				"src": "13.38886",
				"dst": "13.397634,52.529407",
			},
			"Expect `src` to have lattitude and longitude",
			http.StatusBadRequest,
		},
		"malformated `dst` params, no latitude": {
			map[string]string{
				"src": "13.38886,52.517037",
				"dst": "13.397634",
			},
			"Expect 'dst' to have lattitude and longitude",
			http.StatusBadRequest,
		},
		"malformated type for `src` params": {
			map[string]string{
				"src": "la-la-la,52.517037",
				"dst": "13.397634,52.529407",
			},
			"Malformated param type (float64)",
			http.StatusBadRequest,
		},
		"malformated type for `dst` params": {
			map[string]string{
				"src": "13.38886,52.517037",
				"dst": "13.397634,tro-lo-lo",
			},
			"Malformated param type (float64)",
			http.StatusBadRequest,
		},
	}

	for testName, testParams := range testCases {
		t.Run(testName, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodGet, "/routes", nil)
			q := request.URL.Query()
			for key, value := range testParams.params {
				q.Add(key, value)
			}
			request.URL.RawQuery = q.Encode()
			recorder := httptest.NewRecorder()

			GetRoutes(recorder, request)

			responseRecorder := recorder.Result()

			var respMessage views.Response
			err := json.NewDecoder(responseRecorder.Body).Decode(&respMessage)
			if err != nil {
				t.Errorf("%s", err)
			}

			if (respMessage.Body != testParams.message) && (testParams.message != "") {
				t.Errorf("`%v` failed, got %v, expected %v", testName, responseRecorder.Body, testParams.message)
			}

			if responseRecorder.StatusCode != testParams.statusCode {
				t.Errorf("`%v` failed, got %v, expected %v", testName, responseRecorder.StatusCode, testParams.statusCode)
			}
		})
	}
}
