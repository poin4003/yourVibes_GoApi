package query

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type CommentLikeQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func ValidateCommentLikeQueryObject(input interface{}) error {
	query, ok := input.(*CommentLikeQueryObject)
	if !ok {
		return errors.New("input is not CommentLikeQueryObject")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *CommentLikeQueryObject) ToGetCommentLikeQuery(
	commentId uuid.UUID,
) (*commentQuery.GetCommentLikeQuery, error) {
	return &commentQuery.GetCommentLikeQuery{
		CommentId: commentId,
		Limit:     req.Limit,
		Page:      req.Page,
	}, nil
}
