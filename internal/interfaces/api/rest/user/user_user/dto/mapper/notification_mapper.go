package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/response"
)

func MapNotificationToNotificationDto(
	notification *models.Notification,
) *response.NotificationDto {
	return &response.NotificationDto{
		ID:               notification.ID,
		From:             notification.From,
		FromUrl:          notification.FromUrl,
		UserId:           notification.UserId,
		User:             MapUserToUserDtoShortVer(&notification.User),
		NotificationType: notification.NotificationType,
		ContentId:        notification.ContentId,
		Content:          notification.Content,
		Status:           notification.Status,
		CreatedAt:        notification.CreatedAt,
		UpdatedAt:        notification.UpdatedAt,
	}
}
