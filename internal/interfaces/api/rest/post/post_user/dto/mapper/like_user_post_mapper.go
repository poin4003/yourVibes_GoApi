package mapper

import (
	"github.com/google/uuid"
)

func MapToLikeUserPostFromPostIdAndUserId(
	postId uuid.UUID,
	userId uuid.UUID,
) *models.LikeUserPost {
	return &models.LikeUserPost{
		PostId: postId,
		UserId: userId,
	}
}
