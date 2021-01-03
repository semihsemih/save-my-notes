package controllers

import (
	"encoding/json"
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

		if note.Title == "" {
			error.Message = "Title is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if note.Content == "" {
			error.Message = "Content is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if result := c.DB.Create(&note); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusCreated)
		utils.ResponseJSON(w, note)
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

		if result := c.DB.Where("id = ?", note.ID).First(&note); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, note)
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

		if note.Title == "" {
			error.Message = "Title is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if note.Content == "" {
			error.Message = "Content is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if result := c.DB.Model(&note).Updates(models.Note{Title: note.Title, Content: note.Content}); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "Note Updated!")
	}
}

func (c *Controller) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		vars := mux.Vars(r)
		id := vars["id"]

		if result := c.DB.Unscoped().Delete(&models.Note{}, id); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "Note Removed!")
	}
}
