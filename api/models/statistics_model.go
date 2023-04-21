package models

type Statistics struct {
	Value float64 `json:"value"`
	Label string  `json:"label"`
	Type  string  `json:"type"`
}
