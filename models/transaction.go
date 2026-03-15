package models

type Transaction struct {
	ID       int    `json:"id"`
	ItemID   int    `json:"item_id"`
	Type     string `json:"type"` // "IN" atau "OUT"
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}
