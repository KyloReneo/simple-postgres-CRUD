package router

import (
	"github.com/gorilla/mux"

	"github.com/kyloReneo/simple-postgres-CRUD/handlers"
)

//Creating routes
//Router() is exported and used in main.go

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/stocks/{id}", handlers.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock", handlers.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newStock", handlers.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}", handlers.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletestock/{id}", handlers.DeleteStock).Methods("DELETE", "OPTIONS")
	return router
}
