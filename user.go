package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err := decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 7)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash Password")
	}

	user, err := cfg.DB.CreateUser(input.Email, hashPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not Create User")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{user.Id, user.Email})
}

func (cfg *apiConfig) handlerLoginPost(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err := decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return
	}

	dbUsers, err := cfg.DB.GetUsers()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive users")
	}

	id := 0
	email := ""

	for _, dbUser := range dbUsers {
		if input.Email == dbUser.Email {
			err = bcrypt.CompareHashAndPassword(dbUser.Password, []byte(input.Password))

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Incorrect Password")
				return
			}

			id = dbUser.Id
			email = dbUser.Email
		}
	}

	respondWithJSON(w, http.StatusOK, User{id, email})

}
