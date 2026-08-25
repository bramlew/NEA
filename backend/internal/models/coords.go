package models

type Coords struct {
	Lon float64 `json:"lon" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
}
