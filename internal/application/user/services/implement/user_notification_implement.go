package implement

import (
	"context"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/mapper"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"net/http"
)

type sUserNotification struct {
	userRepo         user_repo.IUserRepository
	notificationRepo user_repo.INotificationRepository
}

func NewUserNotificationImplement(
	userRepo user_repo.IUserRepository,
	notificationRepo user_repo.INotificationRepository,
) *sUserNotification {
	return &sUserNotification{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sUserNotification) GetNotificationByUserId(
	ctx context.Context,
	query *user_query.GetManyNotificationQuery,
) (result *user_query.GetManyNotificationQueryResult, err error) {
	result = &user_query.GetManyNotificationQueryResult{}
	// 1. Get notification
	notificationEntities, paging, err := s.notificationRepo.GetMany(ctx, query)
	if err != nil {
		result.Notifications = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
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
	command *user_command.UpdateOneStatusNotificationCommand,
) (result *user_command.UpdateOneStatusNotificationCommandResult, err error) {
	result = &user_command.UpdateOneStatusNotificationCommandResult{}
	notificationUpdateEntity := &user_entity.NotificationUpdate{
		Status: pointer.Ptr(false),
	}

	newNotificationUpdateEntity, err := user_entity.NewNotificationUpdate(notificationUpdateEntity)

	_, err = s.notificationRepo.UpdateOne(ctx, command.NotificationId, newNotificationUpdateEntity)

	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sUserNotification) UpdateManyStatusNotification(
	ctx context.Context,
	command *user_command.UpdateManyStatusNotificationCommand,
) (result *user_command.UpdateManyStatusNotificationCommandResult, err error) {
	result = &user_command.UpdateManyStatusNotificationCommandResult{}
	update_conditions := map[string]interface{}{
		"status":  true,
		"user_id": command.UserId,
	}
	update_data := map[string]interface{}{
		"status": false,
	}

	err = s.notificationRepo.UpdateMany(ctx, update_conditions, update_data)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
