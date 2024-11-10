package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func MapToLikeUserCommentFromCommentIdAndUserId(
	commentId uuid.UUID,
	userId uuid.UUID,
) *models.LikeUserComment {
	return &models.LikeUserComment{
		CommentId: commentId,
		UserId:    userId,
	}
}
