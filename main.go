package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	apiCFG := apiConfig{
		fileserverHits: 0,
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiCFG.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCFG.handlerHits)
	mux.HandleFunc("GET /api/reset", apiCFG.handlerMetricReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidate)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
