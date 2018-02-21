package main

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Logger("Home", r)
	w.Write([]byte("Gorilla!\n"))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	Logger("About", r)
	w.Write([]byte("This is a tux-powered server, using Go, Mux, Redis, and MySQL"))
}