package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/kyloReneo/simple-postgres-CRUD/models"
)

// Create a Response format for dealing with responses
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Creating connection with postgres and handling the error due to the db connection
// Returns a db connection
func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Somthing wrong happed while loading .env file. try again...")
	}

	//Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	//Checking connection that we have made to the database is still alive
	//Handling the db errors
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to the database established successfully.")
	return db
}

// Creates a stock in postgres database
func CreateStock(w http.ResponseWriter, r *http.Request) {

	//Create an empty stock of defined models.stock
	var stock models.Stock

	//Decode the json request
	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body due to the original Error:\n  %v", err)
	}

	//Call the insert stock function for passing the created stock
	insertID := insertStock(stock)

	//Format a response object
	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	//Encode and send the response
	json.NewEncoder(w).Encode(res)
}

// Returns a single stock by its id
func GetStock(w http.ResponseWriter, r *http.Request) {

	//Get the stockid from the request params with th id as the key
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string to int due to the original Error:\n %v", err)
	}

	//Call the getStock function with the requested id
	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock by %v id. Err:\n %v", id, err)
	}

	json.NewEncoder(w).Encode(stock)
}

// Returns all stocks from database
func GetAllStocks(w http.ResponseWriter, r *http.Request) {

	//get all stocks in the db
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stocks due to the original Error:\n %v", err)
	}

	//sending all stocks as response
	json.NewEncoder(w).Encode(stocks)
}

// Updates stock's details in the postgres DB
func UpdateStock(w http.ResponseWriter, r *http.Request) {

	//get the stock id from the request params with the id as key
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body due to the original Error:\n %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// Deletes stock's details in the postgres DB
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records deleted %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// inserts one stock in the DB
func insertStock(stock models.Stock) int64 {

	//Create a postgres db connection
	db := createConnection()

	//Closes the db connection after ending the function
	defer db.Close()

	//Create the sql query
	sqlStatment := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid `
	var id int64

	//Execute the sql query statment
	err := db.QueryRow(sqlStatment, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the insert query due to the original errpr:\n %v", err)
	}

	fmt.Printf("Successfully inserted the record with id  %v", id)
	return id
}

// Gets one stock from the DB by its stockid
func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stock models.Stock

	sqlStatment := `SELECT * FROM stocks WHERE stockid=$1`
	record := db.QueryRow(sqlStatment, id)
	err := record.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row. Error:\n %v", err)
	}

	return stock, err
}

// Gets all stocks from the DB
func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatment := `SELECT * FROM stocks`

	records, err := db.Query(sqlStatment)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer records.Close()

	for records.Next() {
		var stock models.Stock
		err = records.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to sacn the row. Error:\n %v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}

// Updates a stock in the DB by its id
func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatment := `UPDATE stocks SET name=$2 price=$3 company=$4 WHRER stockid=$1`
	newRecord, err := db.Exec(sqlStatment, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. Error:\n %v", err)
	}

	rowsAffected, err := newRecord.RowsAffected()

	if err != nil {
		log.Fatalf("Somthing went wrong while checking the affected rows. Error:\n %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}

// Deletes a single stock in the DB by its id
func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatment := `DELETE FROM stock WHERE stockid=$1`
	res, err := db.Exec(sqlStatment, id)

	if err != nil {
		log.Fatalf("Unable to execute the query due to the original Error:\n %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Somthing went wrong while checking the affected rows. Error:\n %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}
