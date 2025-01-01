package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type UserForAdvertiseResult struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	PhoneNumber *string
	Birthday    *time.Time
	AvatarUrl   string
	CapwallUrl  string
	Privacy     consts.PrivacyLevel
	Biography   string
	PostCount   int
	FriendCount int
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
