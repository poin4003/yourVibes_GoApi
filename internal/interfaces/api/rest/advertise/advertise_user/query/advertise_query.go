package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
)

type AdvertiseQueryObject struct {
	PostId string `form:"post_id,required"`
	Limit  int    `form:"limit,omitempty"`
	Page   int    `form:"page,omitempty"`
}

func ValidateAdvertiseQueryObject(input interface{}) error {
	query, ok := input.(*AdvertiseQueryObject)
	if !ok {
		return fmt.Errorf("validate AdvertiseQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.PostId, validation.Required),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *AdvertiseQueryObject) ToGetManyAdvertiseQuery() (*advertise_query.GetManyAdvertiseQuery, error) {
	var postId uuid.UUID
	if req.PostId != "" {
		parsePostId, err := uuid.Parse(req.PostId)
		if err != nil {
			return nil, err
		}
		postId = parsePostId
	}

	return &advertise_query.GetManyAdvertiseQuery{
		PostId: postId,
		Page:   req.Page,
		Limit:  req.Limit,
	}, nil
}
