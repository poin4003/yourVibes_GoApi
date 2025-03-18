package implement

import (
	"context"
	"github.com/google/uuid"
	notificationCommand "github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/mapper"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/contain"

	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
)

type sNotificationMQ struct {
	notificationRepo repository.INotificationRepository
}

func NewNotification(
	notificationRepo repository.INotificationRepository,
) *sNotificationMQ {
	return &sNotificationMQ{
		notificationRepo: notificationRepo,
	}
}

func (s *sNotificationMQ) HandleBulkNotification(
	ctx context.Context,
	command notificationCommand.NotificationCommand,
	actions []string,
) error {
	// Init notification entity
	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		global.Logger.Error("Failed to parse user ID", zap.Error(err))
		return err
	}

	notif, err := notificationEntity.NewNotification(
		command.From,
		command.FromUrl,
		userID,
		command.NotificationType,
		command.ContentID,
		command.Content,
	)

	var notifications []*notificationEntity.Notification
	if contain.Contains(actions, "db") {
		// Create and get notification from db
		notifications, err = s.notificationRepo.CreateAndGetNotificationsForFriends(ctx, notif)
		if err != nil {
			global.Logger.Error("Failed to create notifications for friends", zap.Error(err))
			return err
		}
	}

	if contain.Contains(actions, "websocket") && global.NotificationSocketHub != nil {
		for _, notification := range notifications {
			socketMsg := mapper.NewNotificationResult(notification)
			if err := global.NotificationSocketHub.SendNotification(notification.UserId.String(), socketMsg); err != nil {
				global.Logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *sNotificationMQ) HandleSingleNotification(
	ctx context.Context,
	command notificationCommand.NotificationCommand,
	actions []string,
) error {
	// Init notification entity
	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		global.Logger.Error("Failed to parse user ID", zap.Error(err))
		return err
	}

	notif, err := notificationEntity.NewNotification(
		command.From,
		command.FromUrl,
		userID,
		command.NotificationType,
		command.ContentID,
		command.Content,
	)

	var notification *notificationEntity.Notification
	if contain.Contains(actions, "db") {
		// Create and get notification
		notification, err = s.notificationRepo.CreateOne(ctx, notif)
		if err != nil {
			global.Logger.Error("Failed to create notification", zap.Error(err))
			return err
		}
	}

	if contain.Contains(actions, "websocket") && global.NotificationSocketHub != nil {
		socketMsg := mapper.NewNotificationResult(notification)
		if err := global.NotificationSocketHub.SendNotification(notification.UserId.String(), socketMsg); err != nil {
			global.Logger.Error("Failed to send notification", zap.Error(err))
		}
	}

	return nil
}
