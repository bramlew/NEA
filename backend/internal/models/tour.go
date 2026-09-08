package models

type Tour struct {
	TotalWeight float64 `json:"totalWeight"`
	Order       []int   `json:"order"`
}
