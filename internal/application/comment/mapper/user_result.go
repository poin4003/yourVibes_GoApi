package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
)

func NewUserResultFromEntity(
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
