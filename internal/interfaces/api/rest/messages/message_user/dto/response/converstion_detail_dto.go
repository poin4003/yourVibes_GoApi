package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ConversationDetailDto struct {
	UserId         uuid.UUID        `json:"user_id"`
	ConversationId uuid.UUID        `json:"conversation_id"`
	User           *UserDto         `json:"user"`
	Conversation   *ConversationDto `json:"conversation"`
	LastMessStatus bool             `json:"last_mess_status"`
	LastMess       *string          `json:"last_mess"`
}

type ConversationDetailWithRoleDto struct {
	UserId           uuid.UUID               `json:"user_id"`
	ConversationId   uuid.UUID               `json:"conversation_id"`
	User             *UserDto                `json:"user"`
	Conversation     *ConversationDto        `json:"conversation"`
	LastMessStatus   bool                    `json:"last_mess_status"`
	LastMess         *string                 `json:"last_mess"`
	ConversationRole consts.ConversationRole `json:"conversation_role"`
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
		LastMess:       conversationDetailResult.LastMess,
	}
}

func ToConversationDetailWithRoleDto(
	conversationDetailResult *common.ConversationDetailResult,
) *ConversationDetailWithRoleDto {
	if conversationDetailResult == nil {
		return nil
	}

	return &ConversationDetailWithRoleDto{
		UserId:           conversationDetailResult.UserId,
		ConversationId:   conversationDetailResult.ConversationId,
		User:             ToUserDto(conversationDetailResult.User),
		Conversation:     ToConversationDto(conversationDetailResult.Conversation),
		LastMessStatus:   conversationDetailResult.LastMessStatus,
		LastMess:         conversationDetailResult.LastMess,
		ConversationRole: conversationDetailResult.ConversationRole,
	}
}

func ToManyConversationDetailDto(conversationDetailResults []*common.ConversationDetailResult) []*ConversationDetailDto {
	if conversationDetailResults == nil {
		return nil
	}

	var dtos []*ConversationDetailDto
	for _, result := range conversationDetailResults {
		dtos = append(dtos, ToConversationDetailDto(result))
	}

	return dtos
}
