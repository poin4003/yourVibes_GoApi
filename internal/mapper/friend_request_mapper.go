package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapToFriendRequestFromUserIdAndFriendId(
	userId uuid.UUID,
	friendId uuid.UUID,
) *model.FriendRequest {
	return &model.FriendRequest{
		UserId:   userId,
		FriendId: friendId,
	}
}
