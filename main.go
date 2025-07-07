package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	fmt.Println("Starting Http Server...")
	var config apiConfig

	mux := http.NewServeMux()

	fileserverHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(fileserverHandler))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets/")))

	//  api/admin endpoints
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", config.handleMetrics)
	mux.HandleFunc("POST /admin/reset", config.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", handleChirp)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// spinning up server
	fmt.Printf("Server running on http://localhost%v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}

func handleChirp(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type:", "application/json")

	type parameters struct {
		Body  string `json:"body"`
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	var response parameters
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&response)
	if err != nil {
		response.Error = fmt.Sprintf("Error decoding json. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response.Error))
		return
	}

	if len(response.Body) > 140 {
		response.Error = "Chirp is too long"
		w.WriteHeader(400)
		w.Write([]byte(response.Error))
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Valid = true

	data, err := json.Marshal(response)
	if err != nil {
		response.Error = "Error marshalling data."
		w.WriteHeader(500)
		return
	}
	w.Write(data)

}
