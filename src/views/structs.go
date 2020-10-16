package views

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
