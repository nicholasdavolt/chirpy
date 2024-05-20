package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
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

	cleaned := cleanChirp(input.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: cleaned,
	})

}

func cleanChirp(input string) string {
	replaceValues := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(input, " ")
	replacement := "****"

	for i, word := range words {
		for _, val := range replaceValues {
			if strings.ToLower(word) == val {
				words[i] = replacement

			}
		}
	}

	return strings.Join(words, " ")

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
