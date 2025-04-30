package command

import (
	"github.com/google/uuid"
)

type CreateConversationDetailCommand struct {
	UserId         uuid.UUID
	ConversationId uuid.UUID
	LastMessStatus bool
	LastMess       string
}

type CreateManyConversationDetailCommand struct {
	UserIds        []uuid.UUID
	ConversationId uuid.UUID
	LastMessStatus bool
	LastMess       string
}
