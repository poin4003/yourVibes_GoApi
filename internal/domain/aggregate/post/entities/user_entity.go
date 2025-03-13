package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	FamilyName string
	Name       string
	AvatarUrl  string
}

func (u *User) ValidateUser() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FamilyName, validation.Required, validation.Length(2, 255)),
		validation.Field(&u.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&u.AvatarUrl, is.URL),
	)
}
