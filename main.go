package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Http Server...")

	mux := http.NewServeMux()

	serveMux := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := serveMux.ListenAndServe()
	if err != nil {
		fmt.Errorf("Error starting server... %v \n", err)
		return
	}

	fmt.Println("Server up and running...")

}
