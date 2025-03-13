package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"gorm.io/gorm"
)

type MessageDto struct {
	ID             uuid.UUID      `json:"id"`
	UserId         uuid.UUID      `json:"user_id"`
	User           *UserDto       `json:"user"`
	ConversationId uuid.UUID      `json:"conversation_id"`
	ParentId       *uuid.UUID     `json:"parent_id"`
	Content        *string        `json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
}

func ToMessageDto(messageResult *common.MessageResult) *MessageDto {
	return &MessageDto{
		ID:             messageResult.ID,
		UserId:         messageResult.UserId,
		User:           ToUserDto(messageResult.User),
		ConversationId: messageResult.ConversationId,
		ParentId:       messageResult.ParentId,
		Content:        messageResult.Content,
		CreatedAt:      messageResult.CreatedAt,
		UpdatedAt:      messageResult.UpdatedAt,
		DeletedAt:      messageResult.DeletedAt,
	}
}
