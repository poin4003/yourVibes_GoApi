package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
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
