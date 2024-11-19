package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `validate:"omitempty,uuid4"`
	FamilyName string    `validate:"required,min=2"`
	Name       string    `validate:"required,min=2"`
	AvatarUrl  string    `validate:"omitempty,url"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
