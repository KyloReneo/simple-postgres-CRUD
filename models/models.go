package models

//Define a simple struct model for records structure in Database
type Stock struct {
	StockId int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
}
