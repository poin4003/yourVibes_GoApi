package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/message/common"
)

type MessageDto struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      UserDto   `json:"user"`
}

func ToMessageDto(messageResult *common.MessageResult) *MessageDto {
	return &MessageDto{
		ID:        messageResult.ID,
		Content:   messageResult.Content,
		CreatedAt: messageResult.CreatedAt,
		User:      *ToUserDto(messageResult.User),
	}
}
