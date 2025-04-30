package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToConversationModel(conversation *entities.CreateConversation) *models.Conversation {
	var conversationDetailModels []*models.ConversationDetail
	for _, conversationDetail := range conversation.ConversationDetail {
		conversationDetailModel := &models.ConversationDetail{
			UserId:           conversationDetail.UserId,
			ConversationId:   conversationDetail.ConversationId,
			LastMessStatus:   false,
			LastMess:         nil,
			ConversationRole: conversationDetail.ConversationRole,
		}
		conversationDetailModels = append(conversationDetailModels, conversationDetailModel)
	}

	var c = &models.Conversation{
		ID:                 conversation.ID,
		Name:               conversation.Name,
		Image:              conversation.Image,
		CreatedAt:          conversation.CreatedAt,
		UpdatedAt:          conversation.UpdatedAt,
		ConversationDetail: conversationDetailModels,
	}

	return c
}

func FromConversationModel(c *models.Conversation) *entities.Conversation {
	var conversation = &entities.Conversation{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	return conversation
}

func FromConversationModelList(conversationModelList []*models.Conversation) []*entities.Conversation {
	conversationEntityList := []*entities.Conversation{}
	for _, conversationModel := range conversationModelList {
		conversationEntityList = append(conversationEntityList, FromConversationModel(conversationModel))
	}
	return conversationEntityList
}
