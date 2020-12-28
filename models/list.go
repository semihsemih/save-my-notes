package models

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
	Notes       []Note `json:"notes"`
}
