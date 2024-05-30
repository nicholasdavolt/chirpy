package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpReceive(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err := decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return

	}

	cleaned, err := validateChirp(input.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not Create Chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve Chirps")
	}

	path := r.PathValue("id")

	id, err := strconv.Atoi(path)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse Id")
	}

	for _, chirp := range dbChirps {
		if chirp.Id == id {
			respondWithJSON(w, http.StatusOK, chirp)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Could not find Id")

}
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve Chirps")
	}

	chirps := []Chirp{}

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:   chirp.Id,
			Body: chirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140

	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return cleanChirp(body), nil

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
