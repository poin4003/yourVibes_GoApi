package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
)

type ConversationDetailDto struct {
	UserId         uuid.UUID        `json:"user_id"`
	ConversationId uuid.UUID        `json:"conversation_id"`
	User           *UserDto         `json:"user"`
	Conversation   *ConversationDto `json:"conversation"`
}

func ToConversationDetailDto(conversationDetailResult *common.ConversationDetailResult) *ConversationDetailDto {
	if conversationDetailResult == nil {
		return nil
	}

	return &ConversationDetailDto{
		UserId:         conversationDetailResult.UserId,
		ConversationId: conversationDetailResult.ConversationId,
		User:           ToUserDto(conversationDetailResult.User),
		Conversation:   ToConversationDto(conversationDetailResult.Conversation),
	}
}
