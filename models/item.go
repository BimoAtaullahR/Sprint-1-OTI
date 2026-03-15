package models

type Item struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Stock        int    `json:"stock"`
	MinimumStock int    `json:"minimum_stock"`
	Unit         string `json:"unit"`
}
