package models

type Response struct {
	Route       []*Leg  `json:"route"`
	TotalWeight float64 `json:"totalWeight"`
}
