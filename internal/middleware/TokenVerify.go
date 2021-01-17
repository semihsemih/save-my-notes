package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"net/http"
	"os"
	"strings"
)

func TokenVerifyMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var errorObject models.Error
	authHeader := r.Header.Get("Authorization")

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) == 2 {
		authToken := bearerToken[1]

		token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if error != nil {
			errorObject.Message = error.Error()
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = error.Error()
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	} else {
		errorObject.Message = "Invalid token"
		utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
}
