package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
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

func NewUserForReportResult(
	user *commentEntity.UserForReport,
) *common.UserForReportResult {
	if user == nil {
		return nil
	}

	return &common.UserForReportResult{
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
		FriendCount: user.FriendCount,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
