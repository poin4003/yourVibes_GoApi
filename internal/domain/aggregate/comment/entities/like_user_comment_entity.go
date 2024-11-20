package entities

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LikeUserComment struct {
	UserId    uuid.UUID `validate:"required,uuid4"`
	CommentId uuid.UUID `validate:"required,uuid4"`
}

func (luc *LikeUserComment) Validate() error {
	validate := validator.New()
	return validate.Struct(luc)
}

func NewLikeUserCommentEntity(
	userId uuid.UUID,
	commentId uuid.UUID,
) (*LikeUserComment, error) {
	fmt.Println(userId)
	fmt.Println(commentId)
	newLikeUserComment := &LikeUserComment{
		UserId:    userId,
		CommentId: commentId,
	}
	if err := newLikeUserComment.Validate(); err != nil {
		return nil, err
	}

	return newLikeUserComment, nil
}
