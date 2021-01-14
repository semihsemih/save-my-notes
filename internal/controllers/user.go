package controllers

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"net/http"
	"os"
	"strconv"
	"time"
)


func (c *Controller) AccountActivation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		vars := mux.Vars(r)
		activationToken := vars["token"]
		var uuid string
		token, err := jwt.Parse(activationToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil {
			errorObject.Message = err.Error()
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			uuid = fmt.Sprintf("%v", claims["uuid"])
		}

		err = c.DB.Exec(
			"UPDATE users SET status = @status, updated_at = @updated_at WHERE uuid = @uuid",
			sql.Named("status", true), sql.Named("updated_at", time.Now()), sql.Named("uuid", uuid),
		).Error
		if err != nil {
			errorObject.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, errorObject)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, "Account activated")
	}
}

func (c *Controller) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		user.ID = uint(id)

		err = c.DB.Where("id = ?", user.ID).First(&user).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}
		err = c.DB.Model(&user).Association("Lists").Find(&user.Lists)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, user)
	}
}
