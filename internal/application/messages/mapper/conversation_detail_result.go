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
		LastMess:       conversationDetail.LastMess,
	}
}

func NewManyConversationDetailResult(
	conversationDetails []*conversationDetailEntity.ConversationDetail,
) []*common.ConversationDetailResult { // ✅ Trả về danh sách thay vì một phần tử
	if conversationDetails == nil {
		return nil
	}

	var results []*common.ConversationDetailResult
	for _, detail := range conversationDetails {
		mappedDetail := &common.ConversationDetailResult{
			UserId:         detail.UserId,
			ConversationId: detail.ConversationId,
			User:           NewMessageUserResultFromEntity(detail.User),
			Conversation:   NewConversationResult(detail.Conversation),
			LastMessStatus: detail.LastMessStatus,
			LastMess:       detail.LastMess,
		}
		results = append(results, mappedDetail) // ✅ Thêm vào danh sách
	}

	return results // ✅ Trả về danh sách đầy đủ
}
