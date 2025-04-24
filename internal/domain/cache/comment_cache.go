package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type (
	ICommentCache interface {
		SetComment(ctx context.Context, comment *entities.Comment)
		GetComment(ctx context.Context, commentID uuid.UUID) *entities.Comment
		DeleteComment(ctx context.Context, commentID uuid.UUID)
		SetPostComment(ctx context.Context, postID uuid.UUID, parentID uuid.UUID, commentIds []uuid.UUID, paging *response.PagingResponse)
		GetPostComment(ctx context.Context, postID uuid.UUID, parentID uuid.UUID, limit, page int) ([]uuid.UUID, *response.PagingResponse)
		DeletePostComment(ctx context.Context, postID uuid.UUID)
		DeleteAllUserComments(ctx context.Context, userID uuid.UUID)
	}
)
