package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromUserModel(
	userModel *models.User,
) *reportEntity.UserForReport {
	if userModel == nil {
		return &reportEntity.UserForReport{}
	}

	var userForReport = &reportEntity.UserForReport{
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
	userForReport.ID = userModel.ID

	return userForReport
}
