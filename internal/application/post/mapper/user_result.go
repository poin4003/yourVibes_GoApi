package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

func NewUserResultFromEntity(
	user *postEntity.User,
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
