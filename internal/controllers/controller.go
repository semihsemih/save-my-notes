package controllers

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Controller struct {
	DB        *gorm.DB
	Validator *validator.Validate
}
