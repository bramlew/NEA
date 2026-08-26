package models

type Request struct {
	Locations []Coords `json:"locations" binding:"required"`
	Weight    float64  `json:"weight"`
}
