package validator

import (
	"fmt"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

type ValidatedPost struct {
	post_entity.Post
	isValidated bool
}

func (vp *ValidatedPost) Valid() bool {
	return vp.isValidated
}

func NewValidatedPost(post *post_entity.Post) (*ValidatedPost, error) {
	if post == nil {
		return nil, fmt.Errorf("post is nil")
	}

	if err := post.ValidatePost(); err != nil {
		return nil, err
	}

	return &ValidatedPost{
		Post:        *post,
		isValidated: true,
	}, nil
}
