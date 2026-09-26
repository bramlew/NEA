package models

type PolylinePayload struct {
	Coordinates      [][]float64 `json:"coordinates"`
	GeometrySimplify bool        `json:"geometry_simplify"`
	Instructions     bool        `json:"instructions"`
	Units            string      `json:"units"`
}
