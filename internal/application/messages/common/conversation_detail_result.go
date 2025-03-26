package common

import "github.com/google/uuid"

type ConversationDetailResult struct {
	UserId         uuid.UUID
	ConversationId uuid.UUID
	User           *UserResult
	Conversation   *ConversationResult
	LastMessStatus bool
	LastMessId     *uuid.UUID
	LastMess       *MessageResult
}
