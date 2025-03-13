package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type CreateMessageCommand struct {
	ConversationId uuid.UUID
	UserId         uuid.UUID
	ParentId       *uuid.UUID
	Content        string
}

type CreateMessageResult struct {
	Message *common.MessageResult
}
