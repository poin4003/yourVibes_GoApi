package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromUserModel(userModel *models.User) *entities.User {
	var user = &entities.User{
		FamilyName: userModel.FamilyName,
		Name:       userModel.Name,
		AvatarUrl:  userModel.AvatarUrl,
	}

	user.ID = userModel.ID

	return user
}
