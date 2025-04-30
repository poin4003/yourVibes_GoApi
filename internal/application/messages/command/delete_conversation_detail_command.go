package command

import "github.com/google/uuid"

type DeleteConversationDetailCommand struct {
	UserId              *uuid.UUID
	ConversationId      *uuid.UUID
	AuthenticatedUserId uuid.UUID
}
