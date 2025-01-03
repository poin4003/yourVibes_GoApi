package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type NewFeedQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func ValidateNewFeedQueryObject(input interface{}) error {
	query, ok := input.(*NewFeedQueryObject)
	if !ok {
		return fmt.Errorf("validate PostLikeQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *NewFeedQueryObject) ToGetNewFeedQuery(
	userId uuid.UUID,
) (*postQuery.GetNewFeedQuery, error) {
	return &postQuery.GetNewFeedQuery{
		UserId: userId,
		Limit:  req.Limit,
		Page:   req.Page,
	}, nil
}
