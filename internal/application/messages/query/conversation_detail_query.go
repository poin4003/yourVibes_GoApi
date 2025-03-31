package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type GetConversationDetailQuery struct {
	ConversationId uuid.UUID
	Limit          int
	Page           int
}

type GetConversationDetailResult struct {
	ConversationDetail []*common.ConversationDetailResult
	PagingResponse     *response.PagingResponse
}
