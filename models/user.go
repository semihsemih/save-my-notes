package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email" gorm:"unique" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6"`
	Status   bool      `json:"status"`
	Lists    []List    `json:"lists"`
}
