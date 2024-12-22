package entities

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type UserForAdvertise struct {
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
