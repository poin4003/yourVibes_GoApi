package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UserWithSettingDto struct {
	ID          uuid.UUID           `json:"id"`
	FamilyName  string              `json:"family_name"`
	Name        string              `json:"name"`
	Email       string              `json:"email"`
	PhoneNumber *string             `json:"phone_number"`
	Birthday    *time.Time          `json:"birthday"`
	AvatarUrl   string              `json:"avatar_url"`
	CapwallUrl  string              `json:"capwall_url"`
	Privacy     consts.PrivacyLevel `json:"privacy"`
	Biography   string              `json:"biography"`
	PostCount   int                 `json:"post_count"`
	FriendCount int                 `json:"friend_count"`
	Status      bool                `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Setting     *SettingDto         `json:"setting"`
}

type SettingDto struct {
	ID        uint            `json:"id"`
	UserId    uuid.UUID       `json:"user_id"`
	Language  consts.Language `json:"language"`
	Status    bool            `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func ToSettingDto(settingResult *common.SettingResult) *SettingDto {
	return &SettingDto{
		ID:        settingResult.ID,
		UserId:    settingResult.UserId,
		Language:  settingResult.Language,
		Status:    settingResult.Status,
		CreatedAt: settingResult.CreatedAt,
		UpdatedAt: settingResult.UpdatedAt,
	}
}

func ToUserWithSettingDto(
	userResult *common.UserWithSettingResult,
) *UserWithSettingDto {
	return &UserWithSettingDto{
		ID:          userResult.ID,
		FamilyName:  userResult.FamilyName,
		Name:        userResult.Name,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
		Birthday:    userResult.Birthday,
		AvatarUrl:   userResult.AvatarUrl,
		CapwallUrl:  userResult.CapwallUrl,
		Privacy:     userResult.Privacy,
		Biography:   userResult.Biography,
		PostCount:   userResult.PostCount,
		FriendCount: userResult.FriendCount,
		Status:      userResult.Status,
		CreatedAt:   userResult.CreatedAt,
		UpdatedAt:   userResult.UpdatedAt,
		Setting:     ToSettingDto(userResult.Setting),
	}
}
