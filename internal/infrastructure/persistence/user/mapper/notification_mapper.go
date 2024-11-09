package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToNotificationModel(notification *user_entity.Notification) *models.Notification {
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

func FromNotificationModel(n *models.Notification) *user_entity.Notification {
	var user = &user_entity.User{
		ID:           n.User.ID,
		FamilyName:   n.User.FamilyName,
		Name:         n.User.Name,
		Email:        n.User.Email,
		Password:     n.User.Password,
		PhoneNumber:  n.User.PhoneNumber,
		Birthday:     n.User.Birthday,
		AvatarUrl:    n.User.AvatarUrl,
		CapwallUrl:   n.User.CapwallUrl,
		Privacy:      n.User.Privacy,
		Biography:    n.User.Biography,
		AuthType:     n.User.AuthType,
		AuthGoogleId: n.User.AuthGoogleId,
		PostCount:    n.User.PostCount,
		FriendCount:  n.User.FriendCount,
		Status:       n.User.Status,
		CreatedAt:    n.User.CreatedAt,
		UpdatedAt:    n.User.UpdatedAt,
	}

	var notification = &user_entity.Notification{
		From:             n.From,
		FromUrl:          n.FromUrl,
		UserId:           n.UserId,
		User:             *user,
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
