package user_dto

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type UserDto struct {
	ID           uuid.UUID           `json:"id"`
	FamilyName   string              `json:"family_name"`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	PhoneNumber  string              `json:"phone_number"`
	Birthday     time.Time           `json:"birthday"`
	AvatarUrl    string              `json:"avatar_url"`
	CapwallUrl   string              `json:"capwall_url"`
	Privacy      consts.PrivacyLevel `json:"validator"`
	Biography    string              `json:"biography"`
	AuthType     string              `json:"auth_type"`
	AuthGoogleId string              `json:"auth_google_id"`
	PostCount    int                 `json:"post_count"`
	FriendCount  int                 `json:"friend_count"`
	Status       bool                `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type UserDtoShortVer struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	AvatarUrl  string    `json:"avatar_url"`
}
