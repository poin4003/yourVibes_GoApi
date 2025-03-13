package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
)

func NewMessageUserResultFromEntity(
	user *entities.User,
) *common.UserResult {
	if user == nil {
		return nil
	}

	return &common.UserResult{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}
