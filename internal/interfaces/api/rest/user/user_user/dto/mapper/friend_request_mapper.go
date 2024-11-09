package mapper

import (
	"github.com/google/uuid"
)

func MapToFriendRequestFromUserIdAndFriendId(
	userId uuid.UUID,
	friendId uuid.UUID,
) *models.FriendRequest {
	return &models.FriendRequest{
		UserId:   userId,
		FriendId: friendId,
	}
}
