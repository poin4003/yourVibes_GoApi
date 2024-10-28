package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapToFriendFromUserIdAndFriendId(
	userId uuid.UUID,
	friendId uuid.UUID,
) *model.Friend {
	return &model.Friend{
		UserId:   userId,
		FriendId: friendId,
	}
}
