package command

import "github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"

type CreateConversationCommand struct {
	Name  string
	Image string
}

type CreateConversationResult struct {
	Conversation *common.ConversationResult
}
