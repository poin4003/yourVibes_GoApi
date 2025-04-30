package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ConversationDetailResult struct {
	UserId           uuid.UUID
	ConversationId   uuid.UUID
	User             *UserResult
	Conversation     *ConversationResult
	LastMessStatus   bool
	LastMess         *string
	ConversationRole consts.ConversationRole
}
