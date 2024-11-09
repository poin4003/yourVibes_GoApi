package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
)

func MapUserToUserDtoWithoutSetting(
	user *models.User,
	friendStatus consts.FriendStatus,
) *response.UserDtoWithoutSetting {
	return &response.UserDtoWithoutSetting{
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
		AuthType:     user.AuthType,
		AuthGoogleId: user.AuthGoogleId,
		PostCount:    user.PostCount,
		Status:       user.Status,
		FriendStatus: friendStatus,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func MapUserToUserDtoShortVer(user *models.User) response.UserDtoShortVer {
	return response.UserDtoShortVer{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}

func MapToUserFromUpdateDto(
	input *request.UpdateUserInput,
) map[string]interface{} {
	updateData := make(map[string]interface{})

	if input.FamilyName != nil {
		updateData["family_name"] = input.FamilyName
	}

	if input.Name != nil {
		updateData["name"] = input.Name
	}

	if input.Email != nil {
		updateData["email"] = input.Email
	}

	if input.PhoneNumber != nil {
		updateData["phone_number"] = input.PhoneNumber
	}

	if input.Birthday != nil {
		updateData["birthday"] = input.Birthday
	}

	if input.Privacy != nil {
		updateData["privacy"] = input.Privacy
	}

	if input.Biography != nil {
		updateData["biography"] = input.Biography
	}

	return updateData
}
