package web

import (
	"eaglebank/internal/validation"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var secretKey = []byte("i-would-not-do-this-in-prod")

func verifyCredentials(_, _ string) bool {
	return true
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err := validation.Get().Struct(req)
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		if !verifyCredentials(req.UserID, req.PasswordHash) {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": req.UserID,
			"exp": now.Add(time.Hour * 24).Unix(),
			"iat": now.Unix(),
		})

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, errors.New("authorization error"))
			return
		}

		resp := LoginResponse{Token: tokenString}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
