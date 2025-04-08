package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/query"
)

type (
	INotificationMQ interface {
		HandleBulkNotification(ctx context.Context, command command.NotificationCommand, actions []string) error
		HandleSingleNotification(ctx context.Context, command command.NotificationCommand, actions []string) error
	}
	INotificationUser interface {
		GetNotificationByUserId(ctx context.Context, query *query.GetManyNotificationQuery) (result *query.GetManyNotificationQueryResult, err error)
		UpdateOneStatusNotification(ctx context.Context, command *command.UpdateOneStatusNotificationCommand) (err error)
		UpdateManyStatusNotification(ctx context.Context, command *command.UpdateManyStatusNotificationCommand) (err error)
	}
)
