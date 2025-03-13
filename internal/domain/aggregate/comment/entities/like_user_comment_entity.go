package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type LikeUserComment struct {
	UserId    uuid.UUID
	CommentId uuid.UUID
}

func (luc *LikeUserComment) ValidateLikeUserComment() error {
	return validation.ValidateStruct(luc,
		validation.Field(&luc.UserId, validation.Required),
		validation.Field(&luc.CommentId, validation.Required),
	)
}

func NewLikeUserCommentEntity(
	userId uuid.UUID,
	commentId uuid.UUID,
) (*LikeUserComment, error) {
	newLikeUserComment := &LikeUserComment{
		UserId:    userId,
		CommentId: commentId,
	}
	if err := newLikeUserComment.ValidateLikeUserComment(); err != nil {
		return nil, err
	}

	return newLikeUserComment, nil
}
