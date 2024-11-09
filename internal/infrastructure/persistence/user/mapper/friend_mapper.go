package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToFriendModel(friend *user_entity.Friend) *models.Friend {
	f := &models.Friend{
		UserId:   friend.UserId,
		FriendId: friend.FriendId,
	}

	return f
}

func FromFriendModel(f *models.Friend) *user_entity.Friend {
	var friend = &user_entity.Friend{
		UserId:   f.UserId,
		FriendId: f.FriendId,
	}

	return friend
}
