package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToNotificationModel(notification *entities.Notification) *models.Notification {
	n := &models.Notification{
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserId:           notification.UserId,
		NotificationType: notification.NotificationType,
		ContentId:        notification.ContentId,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
	n.ID = notification.ID

	return n
}

func FromNotificationModel(n *models.Notification) *entities.Notification {
	var user = &entities.User{
		ID:         n.User.ID,
		FamilyName: n.User.FamilyName,
		Name:       n.User.Name,
		AvatarUrl:  n.User.AvatarUrl,
	}

	var notification = &entities.Notification{
		From:             n.From,
		FromUrl:          n.FromUrl,
		UserId:           n.UserId,
		User:             user,
		NotificationType: n.NotificationType,
		ContentId:        n.ContentId,
		Content:          n.Content,
		Status:           n.Status,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}
	notification.ID = n.ID

	return notification
}
