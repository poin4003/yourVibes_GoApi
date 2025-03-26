package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	conversationDetailEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
)

func NewConversationDetailResult(
	conversationDetail *conversationDetailEntity.ConversationDetail,
) *common.ConversationDetailResult {
	if conversationDetail == nil {
		return nil
	}
	return &common.ConversationDetailResult{
		UserId:         conversationDetail.UserId,
		ConversationId: conversationDetail.ConversationId,
		User:           NewMessageUserResultFromEntity(conversationDetail.User),
		Conversation:   NewConversationResult(conversationDetail.Conversation),
		LastMessStatus: conversationDetail.LastMessStatus,
		LastMessId:     conversationDetail.LastMessId,
		LastMess:       NewMessageResult(conversationDetail.LastMess),
	}
}
