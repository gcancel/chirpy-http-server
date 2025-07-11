package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gcancel/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {

	const port = "8080"
	const filePathDir = "."

	// reading .env parameters
	godotenv.Load()
	db_url := os.Getenv("DB_URL")

	// connecting to database
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error connecting to Database.", err)
	}
	dbQueries := database.New(db)
	config := apiConfig{dbQueries: dbQueries, platform: os.Getenv("PLATFORM")}

	fmt.Println("Starting Http Server...")
	mux := http.NewServeMux()

	fileserverHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathDir)))
	mux.Handle("/app/", config.middlewareMetricsInc(fileserverHandler))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets/")))

	// endpoints
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", config.handleMetrics)
	mux.HandleFunc("POST /admin/reset", config.handleReset)
	mux.HandleFunc("POST /api/chirps", handleChirp)
	mux.HandleFunc("POST /api/users", config.handleUsers)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// spinning up server
	fmt.Printf("Server running on http://localhost%v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
