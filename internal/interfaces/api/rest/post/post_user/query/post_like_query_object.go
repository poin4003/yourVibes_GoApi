package query

import (
	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type PostLikeQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
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
