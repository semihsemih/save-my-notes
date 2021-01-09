package models

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      uint   `json:"user_id"`
	Notes       []Note `json:"notes"`
}
