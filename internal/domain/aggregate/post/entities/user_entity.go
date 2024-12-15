package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type User struct {
	ID         uuid.UUID
	FamilyName string
	Name       string
	AvatarUrl  string
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

func (u *User) ValidateUser() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FamilyName, validation.Required, validation.Length(2, 255)),
		validation.Field(&u.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&u.AvatarUrl, is.URL),
	)
}
