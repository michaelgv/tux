package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("--- Tux Server ---")
	// Create Router
	router := mux.NewRouter()
	// Define Routes
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/about", AboutHandler)
	// Listen On Server
	log.Fatal(http.ListenAndServe(":8000", router))
}
