package mapper

import (
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToFriendModel(friend *userEntity.Friend) *models.Friend {
	f := &models.Friend{
		UserId:   friend.UserId,
		FriendId: friend.FriendId,
	}

	return f
}

func FromFriendModel(f *models.Friend) *userEntity.Friend {
	var friend = &userEntity.Friend{
		UserId:   f.UserId,
		FriendId: f.FriendId,
	}

	return friend
}
