package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type FriendQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func ValidateFriendQueryObject(input interface{}) error {
	query, ok := input.(*FriendQueryObject)
	if !ok {
		return fmt.Errorf("validateFriendQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *FriendQueryObject) ToFriendQuery(
	userId uuid.UUID,
) (*query.FriendQuery, error) {
	return &query.FriendQuery{
		UserId: userId,
		Limit:  req.Limit,
		Page:   req.Page,
	}, nil
}
