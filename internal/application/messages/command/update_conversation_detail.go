package command

import "github.com/google/uuid"

type UpdateOneStatusConversationDetailCommand struct {
	ConversationId uuid.UUID
	UserId         uuid.UUID
}
