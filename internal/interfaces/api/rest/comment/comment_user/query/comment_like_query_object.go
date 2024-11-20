package query

import (
	"github.com/google/uuid"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type CommentLikeQueryObject struct {
	Limit int `form:"limit,omitempty"`
	Page  int `form:"page,omitempty"`
}

func (req *CommentLikeQueryObject) ToGetCommentLikeQuery(
	commentId uuid.UUID,
) (*comment_query.GetCommentLikeQuery, error) {
	return &comment_query.GetCommentLikeQuery{
		CommentId: commentId,
		Limit:     req.Limit,
		Page:      req.Page,
	}, nil
}
