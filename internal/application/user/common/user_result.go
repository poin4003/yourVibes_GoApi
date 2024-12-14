package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type UserWithSettingResult struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	PhoneNumber string
	Birthday    time.Time
	AvatarUrl   string
	CapwallUrl  string
	Privacy     consts.PrivacyLevel
	AuthType    consts.AuthType
	Biography   string
	PostCount   int
	FriendCount int
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Setting     *SettingResult
}

type SettingResult struct {
	ID        uint
	UserId    uuid.UUID
	Language  consts.Language
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserWithoutSettingResult struct {
	ID           uuid.UUID
	FamilyName   string
	Name         string
	Email        string
	PhoneNumber  string
	Birthday     time.Time
	AvatarUrl    string
	CapwallUrl   string
	Privacy      consts.PrivacyLevel
	Biography    string
	PostCount    int
	FriendCount  int
	Status       bool
	FriendStatus consts.FriendStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserForReportResult struct {
	ID          uuid.UUID
	FamilyName  string
	Name        string
	Email       string
	PhoneNumber string
	Birthday    time.Time
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

type UserShortVerResult struct {
	ID         uuid.UUID
	FamilyName string
	Name       string
	AvatarUrl  string
}
