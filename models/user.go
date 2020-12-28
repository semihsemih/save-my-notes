package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Lists    []List `json:"lists"`
}
