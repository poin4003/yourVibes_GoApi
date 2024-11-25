package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type PostLikeQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func ValidatePostLikeQueryObject(input interface{}) error {
	query, ok := input.(*PostLikeQueryObject)
	if !ok {
		return fmt.Errorf("validate PostLikeQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *PostLikeQueryObject) ToGetPostLikeQuery(
	postId uuid.UUID,
) (*post_query.GetPostLikeQuery, error) {
	return &post_query.GetPostLikeQuery{
		PostId: postId,
		Limit:  req.Limit,
		Page:   req.Page,
	}, nil
}
