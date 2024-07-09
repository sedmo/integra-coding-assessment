package models

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// User represents a user in the system
type User struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name" validate:"required,max=50"`
	FirstName  string `json:"first_name" validate:"required,max=50"`
	LastName   string `json:"last_name" validate:"required,max=50"`
	Email      string `json:"email" validate:"required,max=100,email"`
	UserStatus string `json:"user_status" validate:"required,oneof=A I T"`
	Department string `json:"department" validate:"required,max=50"`
}

// Validate validates the User fields.
func (u *User) Validate() error {
	return validate.Struct(u)
}
