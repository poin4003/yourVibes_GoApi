package mapper

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/common"
	entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
)

func NewNotificationResult(
	notification *entity.Notification,
) *common.NotificationResult {
	if notification == nil {
		return nil
	}

	return &common.NotificationResult{
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserID:           fmt.Sprint(notification.UserId),
		NotificationType: notification.NotificationType,
		ContentID:        notification.ContentId,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
}

func NewNotificationResultForInterface(
	notification *userEntity.Notification,
) *common.NotificationResultForInterface {
	if notification == nil {
		return nil
	}

	user := &common.UserShortVerResult{
		ID:         notification.User.ID,
		FamilyName: notification.User.FamilyName,
		Name:       notification.User.Name,
		AvatarUrl:  notification.User.AvatarUrl,
	}

	return &common.NotificationResultForInterface{
		ID:               notification.ID,
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserId:           notification.UserId,
		User:             user,
		NotificationType: notification.NotificationType,
		ContentId:        notification.ContentId,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
}
