package validator

import (
	"fmt"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
)

type ValidatedComment struct {
	comment_entity.Comment
	isValidated bool
}

func (vc *ValidatedComment) Valid() bool {
	return vc.isValidated
}

func NewValidatedComment(comment *comment_entity.Comment) (*ValidatedComment, error) {
	if comment == nil {
		return nil, fmt.Errorf("comment is nil")
	}

	if err := comment.Validate(); err != nil {
		return nil, err
	}

	return &ValidatedComment{
		Comment:     *comment,
		isValidated: true,
	}, nil
}
