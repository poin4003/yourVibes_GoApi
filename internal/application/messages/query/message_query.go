package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type GetMessagesByConversationIdQuery struct {
	ConversationId uuid.UUID
	CreatedAt      time.Time
	SortBy         string
	IsDescending   bool
	Limit          int
	Page           int
}

type GetMessagesByConversationIdResult struct {
	Messages       []*common.MessageResult
	PagingResponse *response.PagingResponse
}
