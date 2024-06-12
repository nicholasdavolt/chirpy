package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) CreateToken(expires int, id int) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Second)),
		Subject:   fmt.Sprintf("%d", id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))

	return tokenString, err

}

func (cfg *apiConfig) CreateRefreshToken(id int) (string, error) {
	byteAmount := 32
	randomBytes := make([]byte, byteAmount)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(randomBytes)

	expiration := time.Now().Add(time.Duration(cfg.RefreshExpiration) * time.Second).UTC()

	cfg.DB.WriteRefreshToken(tokenString, expiration.Format("2006-01-02"), id)

	return tokenString, nil

}

func (cfg *apiConfig) validateExpiration(expiration int) int {

	if expiration > cfg.DefaultExpiration || expiration == 0 {
		return cfg.DefaultExpiration
	}

	return expiration

}

func (cfg *apiConfig) validateToken(tokenString string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, &claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.JwtSecret), nil },
	)

	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()

	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return "", err
	}

	if issuer != "chirpy" {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil

}
