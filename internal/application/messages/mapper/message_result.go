package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
)

func NewMessageResult(
	message *entities.Message,
) *common.MessageResult {
	if message == nil {
		return nil
	}
	return &common.MessageResult{
		ID:             message.ID,
		ConversationId: message.ConversationId,
		Content:        message.Content,
		UserId:         message.UserId,
		User:           NewMessageUserResultFromEntity(message.User),
		ParentId:       message.ParentId,
		ParentContent:  message.ParentContent,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}
}
