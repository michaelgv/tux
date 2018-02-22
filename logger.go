package main

import (
	"net/http"
	"time"
	"log"
)

func Logger(name string, request *http.Request, start time.Time) {
	log.Printf("%s\t%s\t%s\t%s", request.Method,	request.RequestURI, name, time.Since(start), )
}

func GenLogger(name string, start time.Time) {
	log.Printf("%s\t%s\t%s\t%s", "INTERNAL ACTION", name, "Internal Action", time.Since(start))
}