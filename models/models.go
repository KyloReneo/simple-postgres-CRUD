package models

//Define a simple schema for the users table
type Stock struct {
	StockId int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
}
