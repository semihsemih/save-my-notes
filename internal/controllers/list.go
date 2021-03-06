package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (c *Controller) InsertList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list models.List
		var user models.User
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		fmt.Println(bearerToken[1])
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
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
			uuid := fmt.Sprintf("%v", claims["uuid"])
			if err := c.DB.Where("uuid = ?", uuid).First(&user); err.Error != nil {
				errorObject.Message = err.Error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

			list.UserID = user.ID
		}

		json.NewDecoder(r.Body).Decode(&list)

		validate := validator.New()
		err = validate.Struct(list)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				errorObject.Errors = append(errorObject.Errors, fmt.Sprintf("%v", validationError))
			}
			errorObject.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, errorObject)
			return
		}

		err = c.DB.Create(&list).Error
		if err != nil {
			errorObject.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, errorObject)
			return
		}

		utils.ResponseJSON(w, http.StatusCreated, list)
	}
}

func (c *Controller) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list models.List
		var error models.Error
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		list.ID = uint(id)

		err = c.DB.Where("id = ?", list.ID).First(&list).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}
		err = c.DB.Model(&list).Association("Notes").Find(&list.Notes)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, list)
	}
}

func (c *Controller) UpdateList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list models.List
		var error models.Error
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		list.ID = uint(id)
		json.NewDecoder(r.Body).Decode(&list)

		validate := validator.New()
		err = validate.Struct(list)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				error.Errors = append(error.Errors, fmt.Sprintf("%v", validationError))
			}
			error.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		err = c.DB.Model(&list).Updates(models.List{Title: list.Title, Description: list.Description}).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, "List Updated!")
	}
}

func (c *Controller) DeleteList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		vars := mux.Vars(r)
		id := vars["id"]

		err := c.DB.Unscoped().Delete(&models.List{}, id).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, "List Removed!")
	}
}