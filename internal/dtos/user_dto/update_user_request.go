package user_dto

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
	"time"
)

type UpdateUserInput struct {
	FamilyName      *string              `form:"family_name,omitempty"`
	Name            *string              `form:"name,omitempty"`
	Email           *string              `form:"email,omitempty"`
	PhoneNumber     *string              `form:"phone_number,omitempty"`
	Birthday        *time.Time           `form:"birthday,omitempty"`
	AvatarUrl       multipart.FileHeader `form:"avatar_url,omitempty" binding:"omitempty,file"`
	CapwallUrl      multipart.FileHeader `form:"capwall_url,omitempty" binding:"omitempty,file"`
	Privacy         *consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Biography       *string              `form:"biography,omitempty"`
	LanguageSetting *consts.Language     `form:"language_setting,omitempty" binding:"omitempty,language_setting"`
}
