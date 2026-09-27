package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Sudhanshu069/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	apiCfg := apiConfig{
		db:       database.New(db),
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("POST /api/validate_chirp", handleChirp)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminHitsHandler)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	s.ListenAndServe()
}
