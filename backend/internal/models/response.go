package models

type Response struct {
	Locations   []Coords  `json:"locations"`
	Weights     []float64 `json:"weights"`
	TotalWeight float64   `json:"totalWeight"`
}
