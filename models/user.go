package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Status   bool      `json:"status"`
	Lists    []List    `json:"lists"`
}
