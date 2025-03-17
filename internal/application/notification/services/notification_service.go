package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/query"
)

type (
	INotificationMQ interface {
		HandleBulkNotification(ctx context.Context, command command.NotificationCommand) error
		HandleSingleNotification(ctx context.Context, command command.NotificationCommand) error
	}
	INotificationUser interface {
		GetNotificationByUserId(ctx context.Context, query *query.GetManyNotificationQuery) (result *query.GetManyNotificationQueryResult, err error)
		UpdateOneStatusNotification(ctx context.Context, command *command.UpdateOneStatusNotificationCommand) (err error)
		UpdateManyStatusNotification(ctx context.Context, command *command.UpdateManyStatusNotificationCommand) (err error)
	}
)

var (
	localNotificationMQ   INotificationMQ
	localNotificationUser INotificationUser
)

func NotificationMQ() INotificationMQ {
	if localNotificationMQ == nil {
		panic("service_implement localNotificationMQ not found for interface INotificationMQ")
	}

	return localNotificationMQ
}

func NotificationUser() INotificationUser {
	if localNotificationUser == nil {
		panic("service_implement localNotificationUser not found for interface INotificationUser")
	}

	return localNotificationUser
}

func InitNotificationMQ(i INotificationMQ) {
	localNotificationMQ = i
}

func InitNotificationUser(i INotificationUser) {
	localNotificationUser = i
}
