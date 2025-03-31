package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"mime/multipart"
)

type UpdateConversationCommand struct {
	ConversationId *uuid.UUID
	Name           *string
	Image          *multipart.FileHeader
}

type UpdateConversationCommandResult struct {
	Conversation *common.ConversationResult
}
