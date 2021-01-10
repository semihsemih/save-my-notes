package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/semihsemih/save-my-notes/internal/services"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (c *Controller) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)
		validate := validator.New()
		err := validate.Struct(user)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				error.Errors = append(error.Errors, fmt.Sprintf("%v", validationError))
			}
			error.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)
		user.UUID = uuid.NewV4()

		err = c.DB.Create(&user).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		activationToken, _ := utils.GenerateToken(user)
		go services.SendAccountActivationEmail([]string{user.Email}, activationToken)

		user.Password = ""
		utils.ResponseJSON(w, http.StatusCreated, user)
	}
}

func (c *Controller) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)
		validate := validator.New()
		err := validate.Struct(user)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				error.Errors = append(error.Errors, fmt.Sprintf("%v", validationError))
			}
			error.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password
		err = c.DB.Where("email = ?", user.Email).First(&user).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		hashedPassword := user.Password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := utils.GenerateToken(user)
		if err != nil {
			log.Fatal(err)
		}

		jwt.Token = token

		utils.ResponseJSON(w, http.StatusOK, jwt)
	}
}

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
