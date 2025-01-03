package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type CommentQueryObject struct {
	PostId   string `form:"post_id"`
	ParentId string `form:"parent_id,omitempty"`
	Limit    int    `form:"limit,omitempty"`
	Page     int    `form:"page,omitempty"`
}

func ValidateCommentQueryObject(input interface{}) error {
	query, ok := input.(*CommentQueryObject)
	if !ok {
		return fmt.Errorf("input is not CommentQueryObject")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.PostId, validation.Required),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *CommentQueryObject) ToGetManyCommentQuery(
	userId uuid.UUID,
) (*commentQuery.GetManyCommentQuery, error) {
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

	return &commentQuery.GetManyCommentQuery{
		PostId:              postId,
		ParentId:            parentId,
		AuthenticatedUserId: userId,
		Limit:               req.Limit,
		Page:                req.Page,
	}, nil
}
