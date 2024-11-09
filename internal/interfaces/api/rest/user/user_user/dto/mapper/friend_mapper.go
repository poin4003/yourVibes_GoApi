package mapper

import (
	"github.com/google/uuid"
)

func MapToFriendFromUserIdAndFriendId(
	userId uuid.UUID,
	friendId uuid.UUID,
) *models.Friend {
	return &models.Friend{
		UserId:   userId,
		FriendId: friendId,
	}
}
