package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToConversationDetailModel(conversationDetail *entities.ConversationDetail) *models.ConversationDetail {
	cd := &models.ConversationDetail{
		UserId:         conversationDetail.UserId,
		ConversationId: conversationDetail.ConversationId,
	}
	return cd
}

func FromConversationDetailModel(cd *models.ConversationDetail) *entities.ConversationDetail {
	var conversationDetail = &entities.ConversationDetail{
		UserId:         cd.UserId,
		ConversationId: cd.ConversationId,
		User:           FromUserModel(&cd.User),
		Conversation:   FromConversationModel(&cd.Conversation),
	}

	return conversationDetail
}

func FromConversationDetailModelList(conversationDetailModelList []*models.ConversationDetail) []*entities.ConversationDetail {
	conversationDetailEntityList := []*entities.ConversationDetail{}
	for _, conversationDetailModel := range conversationDetailModelList {
		conversationDetailEntityList = append(conversationDetailEntityList, FromConversationDetailModel(conversationDetailModel))
	}
	return conversationDetailEntityList
}
