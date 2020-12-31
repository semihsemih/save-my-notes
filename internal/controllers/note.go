package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (c Controller) InsertNote(db *gorm.DB) http.HandlerFunc {
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

		if result := db.Create(&note); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusCreated)
		utils.ResponseJSON(w, note)
	}
}

func (c Controller) GetNote(db *gorm.DB) http.HandlerFunc {
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

		if result := db.Where("id = ?", note.ID).First(&note); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, note)
	}
}

func (c Controller) UpdateNote(db *gorm.DB) http.HandlerFunc {
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

		if result := db.Model(&note).Updates(models.Note{Title: note.Title, Content: note.Content}); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "Note Updated!")
	}
}

func (c Controller) DeleteNote(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		vars := mux.Vars(r)
		id := vars["id"]

		if result := db.Unscoped().Delete(&models.Note{}, id); result.Error != nil {
			error.Message = result.Error.Error()
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, "Note Removed!")
	}
}
