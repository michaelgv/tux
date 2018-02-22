package main

import (
	"net/http"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Logger("Home", r, time.Now())
	w.Write([]byte("Gorilla!\n"))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	Logger("About", r, time.Now())
	w.Write([]byte("This is a tux-powered server, using Go, Mux, Redis, and MySQL"))
}