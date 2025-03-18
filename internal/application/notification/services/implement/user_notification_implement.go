package implement

import (
	"context"
	notificationCommand "github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	notificationQuery "github.com/poin4003/yourVibes_GoApi/internal/application/notification/query"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"

	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/mapper"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sNotificationUser struct {
	userRepo         userRepo.IUserRepository
	notificationRepo userRepo.INotificationRepository
}

func NewNotificationUserImplement(
	userRepo userRepo.IUserRepository,
	notificationRepo userRepo.INotificationRepository,
) *sNotificationUser {
	return &sNotificationUser{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sNotificationUser) GetNotificationByUserId(
	ctx context.Context,
	query *notificationQuery.GetManyNotificationQuery,
) (result *notificationQuery.GetManyNotificationQueryResult, err error) {
	// 1. Get notification
	notificationEntities, paging, err := s.notificationRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	// 2. Map to result
	var notificationResults []*common.NotificationResultForInterface
	for _, notificationEntity := range notificationEntities {
		notificationResults = append(notificationResults, mapper.NewNotificationResultForInterface(notificationEntity))
	}

	return &notificationQuery.GetManyNotificationQueryResult{
		Notifications:  notificationResults,
		PagingResponse: paging,
	}, nil
}

func (s *sNotificationUser) UpdateOneStatusNotification(
	ctx context.Context,
	command *notificationCommand.UpdateOneStatusNotificationCommand,
) (err error) {
	notificationUpdateEntity := &userEntity.NotificationUpdate{
		Status: pointer.Ptr(false),
	}

	newNotificationUpdateEntity, _ := userEntity.NewNotificationUpdate(notificationUpdateEntity)

	_, err = s.notificationRepo.UpdateOne(ctx, command.NotificationId, newNotificationUpdateEntity)

	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sNotificationUser) UpdateManyStatusNotification(
	ctx context.Context,
	command *notificationCommand.UpdateManyStatusNotificationCommand,
) (err error) {
	updateConditions := map[string]interface{}{
		"status":  true,
		"user_id": command.UserId,
	}
	updateData := map[string]interface{}{
		"status": false,
	}

	err = s.notificationRepo.UpdateMany(ctx, updateConditions, updateData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}
