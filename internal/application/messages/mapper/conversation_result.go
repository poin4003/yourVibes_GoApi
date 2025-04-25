package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	conversationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
)

func NewConversationResult(
	conversation *conversationEntity.Conversation,
) *common.ConversationResult {
	if conversation == nil {
		return nil
	}

	return &common.ConversationResult{
		ID:             conversation.ID,
		Name:           conversation.Name,
		Image:          conversation.Image,
		UserID:         conversation.UserID,
		Avatar:         conversation.Avatar,
		FamilyName:     conversation.FamilyName,
		CreatedAt:      conversation.CreatedAt,
		UpdatedAt:      conversation.UpdatedAt,
		LastMessStatus: conversation.LastMessStatus,
		LastMess:       conversation.LastMess,
	}
}

func NewConversationWithActiveStatusResult(
	conversation *conversationEntity.Conversation,
	activeStatus bool,
) *common.ConversationWithActiveStatusResult {
	if conversation == nil {
		return nil
	}

	return &common.ConversationWithActiveStatusResult{
		ID:             conversation.ID,
		Name:           conversation.Name,
		Image:          conversation.Image,
		UserID:         conversation.UserID,
		Avatar:         conversation.Avatar,
		FamilyName:     conversation.FamilyName,
		CreatedAt:      conversation.CreatedAt,
		UpdatedAt:      conversation.UpdatedAt,
		LastMessStatus: conversation.LastMessStatus,
		LastMess:       conversation.LastMess,
		ActiveStatus:   activeStatus,
	}
}
