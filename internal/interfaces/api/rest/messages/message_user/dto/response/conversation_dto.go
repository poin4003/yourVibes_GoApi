package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type ConversationDto struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToConversationDto(conversationResult *common.ConversationResult) *ConversationDto {
	if conversationResult == nil {
		return nil
	}
	return &ConversationDto{
		ID:        conversationResult.ID,
		Name:      conversationResult.Name,
		Image:     conversationResult.Image,
		CreatedAt: conversationResult.CreatedAt,
		UpdatedAt: conversationResult.UpdatedAt,
	}
}
