package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LikeUserPost struct {
	UserId uuid.UUID `validate:"required,uuid4"`
	PostId uuid.UUID `validate:"required,uuid4"`
}

func (lup *LikeUserPost) Validate() error {
	validate := validator.New()
	return validate.Struct(lup)
}

func NewLikeUserPostEntity(
	userId uuid.UUID,
	postId uuid.UUID,
) (*LikeUserPost, error) {
	newLikeUserPost := &LikeUserPost{
		UserId: userId,
		PostId: postId,
	}
	if err := newLikeUserPost.Validate(); err != nil {
		return nil, err
	}

	return newLikeUserPost, nil
}
