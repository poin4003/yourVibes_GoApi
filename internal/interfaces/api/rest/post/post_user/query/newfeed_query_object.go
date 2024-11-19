package query

import (
	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type NewFeedQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func (req *NewFeedQueryObject) ToGetNewFeedQuery(
	userId uuid.UUID,
) (*post_query.GetNewFeedQuery, error) {
	return &post_query.GetNewFeedQuery{
		UserId: userId,
		Limit:  req.Limit,
		Page:   req.Page,
	}, nil
}
