package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err := decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return

	}

	if len(input.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
