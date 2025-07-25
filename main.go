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
	secretToken    string
	polkaAPIKey    string
}

func main() {

	const port = "8080"
	const filePathDir = "."

	godotenv.Load()
	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error connecting to Database.", err)
	}
	dbQueries := database.New(db)
	config := apiConfig{dbQueries: dbQueries, platform: os.Getenv("PLATFORM"), secretToken: os.Getenv("SECRET_TOKEN"), polkaAPIKey: os.Getenv("POLKA_KEY")}

	fmt.Println("Starting Http Server...")
	mux := http.NewServeMux()

	fileserverHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathDir)))
	mux.Handle("/app/", config.middlewareMetricsInc(fileserverHandler))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets/")))

	// endpoints
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", config.handleMetrics)
	mux.HandleFunc("POST /admin/reset", config.handleReset)
	mux.HandleFunc("POST /api/chirps", config.handleChirpsCreate)
	mux.HandleFunc("GET /api/chirps", config.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.handleGetUserChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", config.handleChirpDelete)
	mux.HandleFunc("POST /api/users", config.handleUsersCreate)
	mux.HandleFunc("PUT /api/users", config.handleUpdateUser)
	mux.HandleFunc("POST /api/login", config.handleLogin)
	mux.HandleFunc("POST /api/refresh", config.handleRefresh)
	mux.HandleFunc("POST /api/revoke", config.handleRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", config.handlePolkaHook)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on http://localhost%v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
