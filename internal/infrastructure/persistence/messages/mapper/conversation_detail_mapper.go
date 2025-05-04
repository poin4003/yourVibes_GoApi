package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToConversationDetailModel(conversationDetail *entities.ConversationDetail) *models.ConversationDetail {
	cd := &models.ConversationDetail{
		UserId:           conversationDetail.UserId,
		ConversationId:   conversationDetail.ConversationId,
		LastMessStatus:   conversationDetail.LastMessStatus,
		LastMess:         conversationDetail.LastMess,
		ConversationRole: conversationDetail.ConversationRole,
	}
	return cd
}

func FromConversationDetailModel(cd *models.ConversationDetail) *entities.ConversationDetail {
	var conversationDetail = &entities.ConversationDetail{
		UserId:           cd.UserId,
		ConversationId:   cd.ConversationId,
		User:             FromUserModel(&cd.User),
		Conversation:     FromConversationModel(&cd.Conversation),
		LastMessStatus:   cd.LastMessStatus,
		LastMess:         cd.LastMess,
		ConversationRole: cd.ConversationRole,
	}

	return conversationDetail
}
