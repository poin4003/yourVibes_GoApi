package command

import "github.com/google/uuid"

type UpdateOneStatusNotificationCommand struct {
	NotificationId uint
}

type UpdateManyStatusNotificationCommand struct {
	UserId uuid.UUID
}
