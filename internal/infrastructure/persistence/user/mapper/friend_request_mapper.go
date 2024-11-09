package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToFriendRequestModel(friendRequest *user_entity.FriendRequest) *models.FriendRequest {
	frq := &models.FriendRequest{
		UserId:   friendRequest.UserId,
		FriendId: friendRequest.FriendId,
	}

	return frq
}

func FromFriendRequestModel(frq *models.FriendRequest) *user_entity.FriendRequest {
	var friendRequest = &user_entity.FriendRequest{
		UserId:   frq.UserId,
		FriendId: frq.FriendId,
	}

	return friendRequest
}
