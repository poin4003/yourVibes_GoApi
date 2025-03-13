package entities

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	Password    string
	PhoneNumber string
	IdentityId  string
	Birthday    time.Time
	Status      bool
	Role        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
