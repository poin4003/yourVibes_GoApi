package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapToLikeUserCommentFromCommentIdAndUserId(
	commentId uuid.UUID,
	userId uuid.UUID,
) *model.LikeUserComment {
	return &model.LikeUserComment{
		CommentId: commentId,
		UserId:    userId,
	}
}
