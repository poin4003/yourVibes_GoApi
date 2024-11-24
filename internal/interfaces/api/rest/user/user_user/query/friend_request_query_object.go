package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type FriendRequestQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func ValidateFriendRequestQueryObject(input interface{}) error {
	query, ok := input.(*FriendRequestQueryObject)
	if !ok {
		return fmt.Errorf("validateFriendRequestQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *FriendRequestQueryObject) ToFriendRequestQuery(
	userId uuid.UUID,
) (*query.FriendRequestQuery, error) {
	return &query.FriendRequestQuery{
		UserId: userId,
		Limit:  req.Limit,
		Page:   req.Page,
	}, nil
}
