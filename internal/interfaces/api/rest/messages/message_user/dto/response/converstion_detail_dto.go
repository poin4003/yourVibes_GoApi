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
	LastMessStatus bool             `json:"last_mess_status"`
	LastMessId     *uuid.UUID       `json:"last_mess_id"`
	LastMess       *MessageDto      `json:"last_mess"`
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
		LastMessStatus: conversationDetailResult.LastMessStatus,
		LastMessId:     conversationDetailResult.LastMessId,
		LastMess:       ToMessageDto(conversationDetailResult.LastMess),
	}
}
