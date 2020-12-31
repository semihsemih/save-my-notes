package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (c Controller) InsertList(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list models.List
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
			fmt.Println(claims)
			id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["id"]), 10, 64)
			if err != nil {
				errorObject.Message = err.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			list.UserID = uint(id)
		}

		json.NewDecoder(r.Body).Decode(&list)

		if list.Title == "" {
			errorObject.Message = "Title is missing."
			utils.RespondWithError(w, http.StatusBadRequest, errorObject)
			return
		}

		if list.Description == "" {
			errorObject.Message = "Description is missing."
			utils.RespondWithError(w, http.StatusBadRequest, errorObject)
			return
		}

		if result := db.Create(&list); result.Error != nil {
			errorObject.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, errorObject)
			return
		}

		w.WriteHeader(http.StatusCreated)
		utils.ResponseJSON(w, list)
	}
}

func (c Controller) GetList(db *gorm.DB) http.HandlerFunc {
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

		if result := db.Where("id = ?", list.ID).First(&list); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}
		err = db.Model(&list).Association("Notes").Find(&list.Notes)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, list)
	}
}

func (c Controller) UpdateList(db *gorm.DB) http.HandlerFunc {
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

		if list.Title == "" {
			error.Message = "Title is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if list.Description == "" {
			error.Message = "Description is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if result := db.Model(&list).Updates(models.List{Title: list.Title, Description: list.Description}); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "List Updated!")
	}
}

func (c Controller) DeleteList(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		vars := mux.Vars(r)
		id := vars["id"]

		if result := db.Unscoped().Delete(&models.List{}, id); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "List Removed!")
	}
}