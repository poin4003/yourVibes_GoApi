package mapper

import (
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	userValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/validator"
)

func NewUserShortVerValidateEntity(
	user *userValidator.ValidatedUser,
) *common.UserShortVerResult {
	return NewUserShortVerEntity(&user.User)
}

func NewUserShortVerEntity(
	user *userEntity.User,
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
	user *userEntity.User,
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
	user *userEntity.User,
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
		PhoneNumber:  nil,
		Birthday:     nil,
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
	user *userValidator.ValidatedUser,
) *common.UserWithSettingResult {
	return NewUserResultFromEntity(&user.User)
}

func NewUserResultFromEntity(
	user *userEntity.User,
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
	setting *userEntity.Setting,
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
