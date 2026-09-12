package models

type MatrixPayload struct {
	Locations    [][]float64 `json:"locations"`
	Destinations []int       `json:"destinations,omitempty"`
	Metrics      []string    `json:"metrics"`
	Sources      []int       `json:"sources,omitempty"`
	Units        string      `json:"units"`
}
