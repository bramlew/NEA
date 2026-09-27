package models

type PolylineResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"`
		} `json:"summary"`
		Geometry string `json:"geometry"`
	} `json:"routes"`
}
