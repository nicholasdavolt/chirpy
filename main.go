package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	database "github.com/nicholasdavolt/chirpy/internal"
)

type apiConfig struct {
	fileserverHits    int
	DB                *database.DB
	JwtSecret         string
	DefaultExpiration int
	RefreshExpiration int
	Polka_Key         string
}

func main() {

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")

	if err != nil {
		log.Printf("DB ERROR %s", err)
	}

	mux := http.NewServeMux()
	apiCFG := apiConfig{
		fileserverHits:    0,
		DB:                db,
		JwtSecret:         jwtSecret,
		DefaultExpiration: 3600,
		RefreshExpiration: 5184000,
		Polka_Key:         polkaKey,
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiCFG.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCFG.handlerHits)
	mux.HandleFunc("GET /api/reset", apiCFG.handlerMetricReset)
	mux.HandleFunc("POST /api/refresh", apiCFG.handlerTokenRefresh)
	mux.HandleFunc("POST /api/revoke", apiCFG.handlerTokenRevoke)
	mux.HandleFunc("POST /api/chirps", apiCFG.handlerChirpReceive)
	mux.HandleFunc("GET /api/chirps", apiCFG.handlerGetChirps)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCFG.handlerDeleteChirp)
	mux.HandleFunc("GET /api/chirps/{id}", apiCFG.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCFG.handlerUserCreate)
	mux.HandleFunc("PUT /api/users", apiCFG.handlerUserPut)
	mux.HandleFunc("POST /api/login", apiCFG.handlerLoginPost)
	mux.HandleFunc("POST /api/polka/webhooks", apiCFG.handlerPolkaPost)

	log.Printf("Serving on port: %s\n", port)

	log.Fatal(srv.ListenAndServe())

}
