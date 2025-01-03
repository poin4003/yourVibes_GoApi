package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
)

func NewNotificationResult(
	notification *userEntity.Notification,
) *common.NotificationResult {
	if notification == nil {
		return nil
	}

	user := &common.UserShortVerResult{
		ID:         notification.User.ID,
		FamilyName: notification.User.FamilyName,
		Name:       notification.User.Name,
		AvatarUrl:  notification.User.AvatarUrl,
	}

	return &common.NotificationResult{
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
