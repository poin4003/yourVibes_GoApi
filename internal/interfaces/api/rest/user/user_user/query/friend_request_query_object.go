package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type FriendRequestQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
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
