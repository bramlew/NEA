package models

type Matrix struct {
	Matrix []*Leg `json:"matrix"`
	Cols   int    `json:"cols"`
}
