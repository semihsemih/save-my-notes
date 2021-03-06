package utils

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/semihsemih/save-my-notes/models"
	"log"
	"net/http"
	"os"
	"time"
)

func RespondWithError(w http.ResponseWriter, status int, err models.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("TOKEN_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": user.UUID,
		"iss":   os.Getenv("APP_NAME"),
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
