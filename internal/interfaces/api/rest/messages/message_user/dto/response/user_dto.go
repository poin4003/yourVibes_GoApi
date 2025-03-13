package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type UserDto struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	AvatarUrl  string    `json:"avatar_url"`
}

func ToUserDto(userResult *common.UserResult) *UserDto {
	if userResult == nil {
		return nil
	}
	return &UserDto{
		ID:         userResult.ID,
		FamilyName: userResult.FamilyName,
		Name:       userResult.Name,
		AvatarUrl:  userResult.AvatarUrl,
	}
}
