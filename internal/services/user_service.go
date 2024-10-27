package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/auth_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/notification_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, in *auth_dto.LoginCredentials) (accessToken string, user *model.User, err error)
		Register(ctx context.Context, in *auth_dto.RegisterCredentials) (resultCode int, err error)
		VerifyEmail(ctx context.Context, email string) (resultCode int, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, userId uuid.UUID) (user *model.User, resultCode int, httpStatusCode int, err error)
		GetManyUsers(ctx context.Context, query *query_object.UserQueryObject) (users []*model.User, resultCode int, httpStatusCode int, response *response.PagingResponse, err error)
		UpdateUser(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}, inAvatarUrl multipart.File, inCapwallUrl multipart.File, languageSetting consts.Language) (user *model.User, resultCode int, httpStatusCode int, err error)
	}
	IUserNotification interface {
		GetNotificationByUserId(ctx context.Context, userId uuid.UUID, query query_object.NotificationQueryObject) (notificationDtos []*notification_dto.NotificationDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error)
		UpdateOneStatusNotification(ctx context.Context, notificationID uint) (resultCode int, httpStatusCode int, err error)
		UpdateManyStatusNotification(ctx context.Context, userId uuid.UUID) (resultCode int, httpStatusCode int, err error)
	}
)

var (
	localUserAuth         IUserAuth
	localUserInfo         IUserInfo
	localUserNotification IUserNotification
)

func UserAuth() IUserAuth {
	if localUserAuth == nil {
		panic("repository_implement localUserLogin not found for interface IUserAuth")
	}

	return localUserAuth
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("repository_implement localUserInfo not found for interface IUserInfo")
	}

	return localUserInfo
}

func UserNotification() IUserNotification {
	if localUserNotification == nil {
		panic("repository_implement localUserNotification not found for interface IUserNotification")
	}

	return localUserNotification
}

func InitUserAuth(i IUserAuth) {
	localUserAuth = i
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}

func InitUserNotification(i IUserNotification) {
	localUserNotification = i
}
