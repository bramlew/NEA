package models

type Request struct {
	Coordinates []Coords `json:"coordinates" binding:"required"`
	Weight      float64  `json:"weight"`
}
