package models

import (
	"github.com/go-playground/validator/v10"
)

// User example
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=4,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (u *User) Validate() error {
	validate := validator.New()

	return validate.Struct(u)
}