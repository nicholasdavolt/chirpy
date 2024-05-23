package main

import (
	"log"
	"net/http"

	database "github.com/nicholasdavolt/chirpy/internal"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")

	if err != nil {
		log.Printf("DB ERROR %s", err)
	}

	mux := http.NewServeMux()
	apiCFG := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiCFG.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCFG.handlerHits)
	mux.HandleFunc("GET /api/reset", apiCFG.handlerMetricReset)
	mux.HandleFunc("POST /api/chirps", apiCFG.handlerChirpReceive)
	mux.HandleFunc("GET /api/chirps", apiCFG.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCFG.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCFG.handlerUserCreate)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
