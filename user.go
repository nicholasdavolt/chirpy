package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

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
		Password           string `json:"password"`
		Email              string `json:"email"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err := decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode input: %v", err))
		return
	}

	dbUsers, err := cfg.DB.GetUsers()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive users")
	}

	id := 0
	email := ""
	expiresInSeconds := 0

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

	expiresInSeconds = cfg.validateExpiration(input.Expires_in_seconds)

	tokenString, err := cfg.CreateToken(expiresInSeconds, id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't generate token string: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, UserLogin{id, email, tokenString})

}

func (cfg *apiConfig) handlerUserPut(w http.ResponseWriter, r *http.Request) {
	type inputs struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	authHeader := r.Header.Get("Authorization")

	splitHeader := strings.Split(authHeader, " ")

	token := splitHeader[1]

	userIDString, err := cfg.validateToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	input := inputs{}
	err = decoder.Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode input: %v", err))
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 7)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash Password")
	}

	user, err := cfg.DB.UpdateUser(userIDString, input.Email, hashPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not write edited user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{user.Id, user.Email})
}
