package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"net/http"
	"strconv"
)

func (c *Controller) InsertNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.Note
		var error models.Error
		json.NewDecoder(r.Body).Decode(&note)

		validate := validator.New()
		err := validate.Struct(note)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				error.Errors = append(error.Errors, fmt.Sprintf("%v", validationError))
			}
			error.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		err = c.DB.Create(&note).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		utils.ResponseJSON(w, http.StatusCreated, note)
	}
}

func (c *Controller) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.Note
		var error models.Error
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		note.ID = uint(id)

		err = c.DB.Where("id = ?", note.ID).First(&note).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, note)
	}
}

func (c *Controller) UpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.Note
		var error models.Error
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		note.ID = uint(id)
		json.NewDecoder(r.Body).Decode(&note)

		validate := validator.New()
		err = validate.Struct(note)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				error.Errors = append(error.Errors, fmt.Sprintf("%v", validationError))
			}
			error.Message = "Invalid Request Payload"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		err = c.DB.Model(&note).Updates(models.Note{Title: note.Title, Content: note.Content}).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, "Note Updated!")
	}
}

func (c *Controller) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		vars := mux.Vars(r)
		id := vars["id"]

		err := c.DB.Unscoped().Delete(&models.Note{}, id).Error
		if err != nil {
			error.Message = err.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		utils.ResponseJSON(w, http.StatusOK, "Note Removed!")
	}
}
