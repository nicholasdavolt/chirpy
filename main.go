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
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCFG.handlerHits)
	mux.HandleFunc("/reset", apiCFG.handlerMetricReset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
