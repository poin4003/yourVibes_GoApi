package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type CreateConversationDetailCommand struct {
	UserId         uuid.UUID
	ConversationId uuid.UUID
}

type CreateConversationDetailResult struct {
	ConversationDetail *common.ConversationDetailResult
}
