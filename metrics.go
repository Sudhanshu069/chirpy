package main

import (
	"fmt"
	"net/http"
)

// func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
// 	hits := cfg.fileserverHits.Load()
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	fmt.Fprintf(w, "Hits: %d\n", hits)
// }

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not Allowed")
		return
	}

	cfg.fileserverHits.Swap(0)

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) adminHitsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
