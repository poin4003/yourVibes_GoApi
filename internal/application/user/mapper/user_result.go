package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	user_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
	"time"
)

func NewUserShortVerValidateEntity(
	user *user_validator.ValidatedUser,
) *common.UserShortVerResult {
	return NewUserShortVerEntity(&user.User)
}

func NewUserShortVerEntity(
	user *user_entity.User,
) *common.UserShortVerResult {
	if user == nil {
		return nil
	}

	return &common.UserShortVerResult{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}

func NewUserResultWithoutSettingEntity(
	user *user_entity.User,
	friendStatus consts.FriendStatus,
) *common.UserWithoutSettingResult {
	if user == nil || friendStatus == "" {
		return nil
	}

	return &common.UserWithoutSettingResult{
		ID:           user.ID,
		FamilyName:   user.FamilyName,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Birthday:     user.Birthday,
		AvatarUrl:    user.AvatarUrl,
		CapwallUrl:   user.CapwallUrl,
		Privacy:      user.Privacy,
		Biography:    user.Biography,
		PostCount:    user.PostCount,
		FriendCount:  user.FriendCount,
		Status:       user.Status,
		FriendStatus: friendStatus,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func NewUserResultWithoutPrivateInfo(
	user *user_entity.User,
	friendStatus consts.FriendStatus,
) *common.UserWithoutSettingResult {
	if user == nil || friendStatus == "" {
		return nil
	}

	return &common.UserWithoutSettingResult{
		ID:           user.ID,
		FamilyName:   user.FamilyName,
		Name:         user.Name,
		Email:        "",
		PhoneNumber:  "",
		Birthday:     time.Time{},
		AvatarUrl:    user.AvatarUrl,
		CapwallUrl:   user.CapwallUrl,
		Privacy:      user.Privacy,
		Biography:    "",
		PostCount:    0,
		FriendCount:  0,
		Status:       user.Status,
		FriendStatus: friendStatus,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}

func NewUserResultFromValidateEntity(
	user *user_validator.ValidatedUser,
) *common.UserWithSettingResult {
	return NewUserResultFromEntity(&user.User)
}

func NewUserResultFromEntity(
	user *user_entity.User,
) *common.UserWithSettingResult {
	if user == nil {
		return nil
	}

	return &common.UserWithSettingResult{
		ID:          user.ID,
		FamilyName:  user.FamilyName,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Birthday:    user.Birthday,
		AvatarUrl:   user.AvatarUrl,
		CapwallUrl:  user.CapwallUrl,
		Privacy:     user.Privacy,
		Biography:   user.Biography,
		PostCount:   user.PostCount,
		FriendCount: user.FriendCount,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Setting:     NewSettingResultFromEntity(user.Setting),
	}
}

func NewSettingResultFromEntity(
	setting *user_entity.Setting,
) *common.SettingResult {
	if setting == nil {
		return nil
	}

	return &common.SettingResult{
		ID:        setting.ID,
		UserId:    setting.UserId,
		Language:  setting.Language,
		Status:    setting.Status,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
