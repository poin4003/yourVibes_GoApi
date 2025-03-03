package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetCommentLikeQuery struct {
	CommentId uuid.UUID
	Limit     int
	Page      int
}

type GetCommentLikeResult struct {
	Users          []*common.UserResult
	PagingResponse *response.PagingResponse
}
