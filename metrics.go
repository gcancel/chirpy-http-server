package main

import (
	"context"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf(`
<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>

</html>`, hits)
	w.Write([]byte(message))
}
func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not Allowed", fmt.Errorf("forbidden action"))
		return
	}
	ctx := context.Background()
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteAllUsers(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error clearing users table.", err)
	}
	err = cfg.dbQueries.DeleteAllTokens(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting tokens", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset!\nUsers Database cleared!\nJWT Tokens Cleared!"))
}
