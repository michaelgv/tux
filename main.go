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
	// Define "API" Routes for Authentication
	router.HandleFunc("/api/v1/users/create", ApiCreateUserRoute).Methods("POST")
	router.HandleFunc("/api/v1/users/get", ApiGetUserRoute).Methods("POST")
	router.HandleFunc("/api/v1/users/login", ApiAuthenticateUserRoute).Methods("POST")
	router.HandleFunc("/api/v1/users/changepassword", ApiAuthenticationChangePasswordRoute).Methods("POST")
	router.HandleFunc("/api/v1/users/list", ApiListUsersSafelyRoute).Methods("GET", "POST")
	// Listen On Server
	log.Fatal(http.ListenAndServe(":8000", router))
}
