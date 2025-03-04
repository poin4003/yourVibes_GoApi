package message

import "github.com/google/uuid"

// CreateMessageRequest represents the request body for creating a new message.
type CreateMessageRequest struct {
	ConversationID uuid.UUID `json:"conversationId"`
	SenderID       uuid.UUID `json:"senderId"`
	Content        string    `json:"content"`
}
