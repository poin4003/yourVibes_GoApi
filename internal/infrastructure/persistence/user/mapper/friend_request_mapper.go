package mapper

import (
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToFriendRequestModel(friendRequest *userEntity.FriendRequest) *models.FriendRequest {
	frq := &models.FriendRequest{
		UserId:   friendRequest.UserId,
		FriendId: friendRequest.FriendId,
	}

	return frq
}

func FromFriendRequestModel(frq *models.FriendRequest) *userEntity.FriendRequest {
	var friendRequest = &userEntity.FriendRequest{
		UserId:   frq.UserId,
		FriendId: frq.FriendId,
	}

	return friendRequest
}
