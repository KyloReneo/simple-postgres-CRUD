package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/kyloReneo/simple-postgres-CRUD/routers"

)

func main() {
	r := router.Router()
	fmt.Println("Starting Server on the port 8081...")

	log.Fatal(http.ListenAndServe(":8081", r))
}
