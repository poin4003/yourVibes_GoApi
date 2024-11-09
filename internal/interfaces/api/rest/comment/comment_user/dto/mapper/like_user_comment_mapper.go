package mapper

import (
	"github.com/google/uuid"
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
