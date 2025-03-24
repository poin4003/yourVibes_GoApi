package mapper

import (
	advetiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromUserModel(
	userModel *models.User,
) *advetiseEntity.UserForAdvertise {
	if userModel == nil {
		return &advetiseEntity.UserForAdvertise{}
	}

	var userForAdvertise = &advetiseEntity.UserForAdvertise{
		FamilyName:  userModel.FamilyName,
		Name:        userModel.Name,
		Email:       userModel.Email,
		PhoneNumber: userModel.PhoneNumber,
		Birthday:    userModel.Birthday,
		AvatarUrl:   userModel.AvatarUrl,
		CapwallUrl:  userModel.CapwallUrl,
		Privacy:     userModel.Privacy,
		Biography:   userModel.Biography,
		PostCount:   userModel.PostCount,
		FriendCount: userModel.FriendCount,
		Status:      userModel.Status,
		CreatedAt:   userModel.CreatedAt,
		UpdatedAt:   userModel.UpdatedAt,
	}
	userForAdvertise.ID = userModel.ID

	return userForAdvertise
}
