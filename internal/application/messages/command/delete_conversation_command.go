package command

import "github.com/google/uuid"

type DeleteConversationCommand struct {
	ConversationId *uuid.UUID
	UserId         *uuid.UUID
}
