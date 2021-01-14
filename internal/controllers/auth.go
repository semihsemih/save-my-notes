package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"github.com/semihsemih/save-my-notes/internal/services"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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
