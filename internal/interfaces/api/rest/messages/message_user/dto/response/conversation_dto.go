package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type ConversationDto struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Image          string     `json:"image"`
	Avatar         string     `json:"avatar,omitempty"`
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	FamilyName     string     `json:"family_name,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastMess       *string    `json:"last_message"`
	LastMessStatus bool       `json:"last_message_status"`
}

func ToConversationDto(conversationResult *common.ConversationResult) *ConversationDto {
	if conversationResult == nil {
		return nil
	}
	return &ConversationDto{
		ID:             conversationResult.ID,
		Name:           conversationResult.Name,
		Image:          conversationResult.Image,
		UserID:         conversationResult.UserID,
		Avatar:         conversationResult.Avatar,
		FamilyName:     conversationResult.FamilyName,
		CreatedAt:      conversationResult.CreatedAt,
		UpdatedAt:      conversationResult.UpdatedAt,
		LastMess:       conversationResult.LastMess,
		LastMessStatus: conversationResult.LastMessStatus,
	}
}
