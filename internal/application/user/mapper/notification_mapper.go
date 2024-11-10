package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
)

func NewNotificationResult(
	notification *user_entity.Notification,
) *common.NotificationResult {
	if notification == nil {
		return nil
	}

	return &common.NotificationResult{
		ID:               notification.ID,
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserId:           notification.UserId,
		User:             *NewUserShortVerEntity(&notification.User),
		NotificationType: notification.NotificationType,
		ContentId:        notification.ContentId,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
}
