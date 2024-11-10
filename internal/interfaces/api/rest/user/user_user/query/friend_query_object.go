package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type FriendQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
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
