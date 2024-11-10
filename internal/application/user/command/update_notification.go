package command

import "github.com/google/uuid"

type UpdateOneStatusNotificationCommand struct {
	NotificationId uint
}

type UpdateManyStatusNotificationCommand struct {
	UserId uuid.UUID
}

type UpdateOneStatusNotificationCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}

type UpdateManyStatusNotificationCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
