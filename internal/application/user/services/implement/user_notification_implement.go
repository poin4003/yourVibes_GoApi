package implement

import (
	"context"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	userRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"net/http"
)

type sUserNotification struct {
	userRepo         userRepo.IUserRepository
	notificationRepo userRepo.INotificationRepository
}

func NewUserNotificationImplement(
	userRepo userRepo.IUserRepository,
	notificationRepo userRepo.INotificationRepository,
) *sUserNotification {
	return &sUserNotification{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sUserNotification) GetNotificationByUserId(
	ctx context.Context,
	query *userQuery.GetManyNotificationQuery,
) (result *userQuery.GetManyNotificationQueryResult, err error) {
	result = &userQuery.GetManyNotificationQueryResult{}
	result.Notifications = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Get notification
	notificationEntities, paging, err := s.notificationRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	// 2. Map to result
	var notificationResults []*common.NotificationResult
	for _, notificationEntity := range notificationEntities {
		notificationResults = append(notificationResults, mapper.NewNotificationResult(notificationEntity))
	}

	result.Notifications = notificationResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sUserNotification) UpdateOneStatusNotification(
	ctx context.Context,
	command *userCommand.UpdateOneStatusNotificationCommand,
) (result *userCommand.UpdateOneStatusNotificationCommandResult, err error) {
	result = &userCommand.UpdateOneStatusNotificationCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	notificationUpdateEntity := &userEntity.NotificationUpdate{
		Status: pointer.Ptr(false),
	}

	newNotificationUpdateEntity, err := userEntity.NewNotificationUpdate(notificationUpdateEntity)

	_, err = s.notificationRepo.UpdateOne(ctx, command.NotificationId, newNotificationUpdateEntity)

	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserNotification) UpdateManyStatusNotification(
	ctx context.Context,
	command *userCommand.UpdateManyStatusNotificationCommand,
) (result *userCommand.UpdateManyStatusNotificationCommandResult, err error) {
	result = &userCommand.UpdateManyStatusNotificationCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	updateConditions := map[string]interface{}{
		"status":  true,
		"user_id": command.UserId,
	}
	updateData := map[string]interface{}{
		"status": false,
	}

	err = s.notificationRepo.UpdateMany(ctx, updateConditions, updateData)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
