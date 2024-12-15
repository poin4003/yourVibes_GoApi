package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type User struct {
	ID         uuid.UUID `validate:"omitempty,uuid4"`
	FamilyName string    `validate:"required,min=2"`
	Name       string    `validate:"required,min=2"`
	AvatarUrl  string    `validate:"omitempty,url"`
}

type UserForReport struct {
	ID           uuid.UUID
	FamilyName   string
	Name         string
	Email        string
	Password     string
	PhoneNumber  string
	Birthday     time.Time
	AvatarUrl    string
	CapwallUrl   string
	Privacy      consts.PrivacyLevel
	Biography    string
	AuthType     consts.AuthType
	AuthGoogleId string
	PostCount    int
	FriendCount  int
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
