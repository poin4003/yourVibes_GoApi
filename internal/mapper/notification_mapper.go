package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/notification_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapNotificationToNotificationDto(
	notification *model.Notification,
) *notification_dto.NotificationDto {
	return &notification_dto.NotificationDto{
		ID:               notification.ID,
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserId:           notification.UserId,
		User:             MapUserToUserDtoShortVer(&notification.User),
		NotificationType: notification.NotificationType,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
}
