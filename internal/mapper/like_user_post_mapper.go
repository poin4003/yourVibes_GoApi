package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapToLikeUserPostFromPostIdAndUserId(
	postId uuid.UUID,
	userId uuid.UUID,
) *model.LikeUserPost {
	return &model.LikeUserPost{
		PostId: postId,
		UserId: userId,
	}
}
