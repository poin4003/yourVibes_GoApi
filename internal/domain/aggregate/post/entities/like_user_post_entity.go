package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type LikeUserPost struct {
	UserId uuid.UUID
	PostId uuid.UUID
}

func (lup *LikeUserPost) ValidateLikeUserPost() error {
	return validation.ValidateStruct(lup,
		validation.Field(&lup.UserId, validation.Required),
		validation.Field(&lup.PostId, validation.Required),
	)
}

func NewLikeUserPostEntity(
	userId uuid.UUID,
	postId uuid.UUID,
) (*LikeUserPost, error) {
	newLikeUserPost := &LikeUserPost{
		UserId: userId,
		PostId: postId,
	}
	if err := newLikeUserPost.ValidateLikeUserPost(); err != nil {
		return nil, err
	}

	return newLikeUserPost, nil
}
