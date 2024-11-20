package query

import (
	"github.com/google/uuid"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type CommentQueryObject struct {
	PostId   string `form:"post_id" binding:"required"`
	ParentId string `form:"parent_id,omitempty"`
	Limit    int    `form:"limit,omitempty"`
	Page     int    `form:"page,omitempty"`
}

func (req *CommentQueryObject) ToGetManyCommentQuery(
	userId uuid.UUID,
) (*comment_query.GetManyCommentQuery, error) {
	var postId uuid.UUID
	var parentId uuid.UUID

	if req.PostId != "" {
		parsePostId, err := uuid.Parse(req.PostId)
		if err != nil {
			return nil, err
		}
		postId = parsePostId
	}

	if req.ParentId == "" {
		parentId = uuid.Nil
	} else {
		parseParentId, err := uuid.Parse(req.ParentId)
		if err != nil {
			return nil, err
		}
		parentId = parseParentId
	}

	return &comment_query.GetManyCommentQuery{
		PostId:              postId,
		ParentId:            parentId,
		AuthenticatedUserId: userId,
		Limit:               req.Limit,
		Page:                req.Page,
	}, nil
}
