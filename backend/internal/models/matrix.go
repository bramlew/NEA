package models

type Matrix struct {
	Matrix []*Route `json:"matrix"`
	Cols   int      `json:"cols"`
}
