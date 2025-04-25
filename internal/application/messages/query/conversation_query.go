package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type GetOneConversationQuery struct {
	ConversationId uuid.UUID
}

type GetManyConversationQuery struct {
	Name         string
	CreatedAt    time.Time
	SortBy       string
	IsDescending bool
	Limit        int
	Page         int
}

type GetManyConversationQueryResult struct {
	Conversation   []*common.ConversationWithActiveStatusResult
	PagingResponse *response.PagingResponse
}
