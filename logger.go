package main

import (
	"net/http"
	"time"
	"log"
)

func Logger(name string, request *http.Request) {
	start := time.Now()
	log.Printf("%s\t%s\t%s\t%s", request.Method,	request.RequestURI, name, time.Since(start), )
}
