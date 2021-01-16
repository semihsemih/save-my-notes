package controllers

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/utils"
	"github.com/semihsemih/save-my-notes/models"
	"net/http"
	"strconv"
)

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
