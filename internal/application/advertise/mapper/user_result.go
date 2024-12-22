package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	advertise_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

func NewUserForAdvertiseResult(
	user *advertise_entity.UserForAdvertise,
) *common.UserForAdvertiseResult {
	if user == nil {
		return nil
	}

	return &common.UserForAdvertiseResult{
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
