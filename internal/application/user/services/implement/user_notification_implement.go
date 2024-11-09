package implement

import (
	"context"
	"github.com/google/uuid"
	user_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/mapper"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
	userId uuid.UUID,
	query query.NotificationQueryObject,
) (notificationDtos []*response2.NotificationDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error) {
	notificationModels, paging, err := s.notificationRepo.GetManyNotification(ctx, userId, &query)
	if err != nil {
		return nil, nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	for _, notification := range notificationModels {
		notificationDto := mapper.MapNotificationToNotificationDto(notification)
		notificationDtos = append(notificationDtos, notificationDto)
	}

	return notificationDtos, paging, response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserNotification) UpdateOneStatusNotification(
	ctx context.Context,
	notificationID uint,
) (resultCode int, httpStatusCode int, err error) {
	_, err = s.notificationRepo.UpdateOneNotification(ctx, notificationID, map[string]interface{}{
		"status": false,
	})

	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, err
	}

	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sUserNotification) UpdateManyStatusNotification(
	ctx context.Context,
	userId uuid.UUID,
) (resultCode int, httpStatusCode int, err error) {
	update_conditions := map[string]interface{}{
		"status":  true,
		"user_id": userId,
	}
	update_data := map[string]interface{}{
		"status": false,
	}

	err = s.notificationRepo.UpdateManyNotification(ctx, update_conditions, update_data)
	if err != nil {
		return response.ErrServerFailed, http.StatusInternalServerError, err
	}

	return response.ErrCodeSuccess, http.StatusOK, nil
}
