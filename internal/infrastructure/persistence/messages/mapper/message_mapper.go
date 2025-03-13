package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToMessageModel(message *entities.Message) *models.Message {
	var m = &models.Message{
		ID:             message.ID,
		Content:        message.Content,
		ConversationId: message.ConversationId,
		UserId:         message.UserId,
		ParentId:       message.ParentId,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}
	return m
}
func ToUserEntity(
	user *models.User,
) *entities.User {
	if user == nil {
		return nil
	}

	return &entities.User{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}

func FromMessageModel(m *models.Message) *entities.Message {
	if m == nil {
		return nil
	}

	return &entities.Message{
		ID:             m.ID,
		UserId:         m.UserId,
		User:           ToUserEntity(&m.User),
		ConversationId: m.ConversationId,
		ParentId:       m.ParentId,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
