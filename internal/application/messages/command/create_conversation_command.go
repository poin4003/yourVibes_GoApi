package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type CreateConversationCommand struct {
	Name    string
	Image   string
	UserIds []uuid.UUID
	OwnerId uuid.UUID
}

type CreateConversationResult struct {
	Conversation *common.ConversationResult
}
