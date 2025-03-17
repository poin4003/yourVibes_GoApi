package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/common"
)

type UserShortVerDto struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	AvatarUrl  string    `json:"avatar_url"`
}

func ToUserShortVerDto(
	userResult *common.UserShortVerResult,
) *UserShortVerDto {
	return &UserShortVerDto{
		ID:         userResult.ID,
		FamilyName: userResult.FamilyName,
		Name:       userResult.Name,
		AvatarUrl:  userResult.AvatarUrl,
	}
}
