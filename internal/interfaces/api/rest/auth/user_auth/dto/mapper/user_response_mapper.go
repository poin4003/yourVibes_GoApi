package mapper

import (
	user_common "github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/response"
)

func ToUserWithSettingResponse(user *user_common.UserWithSettingResult) *response.UserResponse {
	return &response.UserResponse{
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
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Setting:     ToSettingResponse(user.Setting),
	}
}

func ToSettingResponse(setting *user_common.SettingResult) response.SettingResponse {
	return response.SettingResponse{
		ID:        setting.ID,
		UserId:    setting.UserId,
		Language:  setting.Language,
		Status:    setting.Status,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
