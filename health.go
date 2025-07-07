package main

import (
	"fmt"
	"net/http"
)

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Starting health check.\n")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}
