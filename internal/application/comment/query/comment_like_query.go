package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type GetCommentLikeQuery struct {
	CommentId uuid.UUID
	Limit     int
	Page      int
}

type CheckUserLikeManyCommentQuery struct {
	CommentIds          []uuid.UUID
	AuthenticatedUserId uuid.UUID
}

type GetCommentLikeResult struct {
	Users          []*common.UserResult
	PagingResponse *response.PagingResponse
}
