package mapper

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
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
