package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type GetManyCommentQuery struct {
	AuthenticatedUserId uuid.UUID
	PostId              uuid.UUID
	ParentId            uuid.UUID
	Limit               int
	Page                int
}

type GetManyCommentsResult struct {
	Comments       []*common.CommentResultWithLiked
	PagingResponse *response.PagingResponse
}
