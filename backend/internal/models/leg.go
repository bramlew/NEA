package models

type Leg struct {
	Distance float64 `json:"distance"`
	Origin   Coords  `json:"origin"`
	Dest     Coords  `json:"dest"`
	IsRoad   bool    `json:"isRoad"`
	Polyline string  `json:"polyline,omitempty"`
}
