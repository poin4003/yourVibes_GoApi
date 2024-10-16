package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/auth_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"mime/multipart"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, in *auth_dto.LoginCredentials) (accessToken string, user *model.User, err error)
		Register(ctx context.Context, in *auth_dto.RegisterCredentials) (resultCode int, err error)
		VerifyEmail(ctx context.Context, email string) (resultCode int, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, userId uuid.UUID) (user *model.User, resultCode int, err error)
		GetManyUsers(ctx context.Context, query *query_object.UserQueryObject) (users []*model.User, resultCode int, err error)
		UpdateUser(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}, inAvatarUrl multipart.File, inCapwallUrl multipart.File, languageSetting consts.Language) (user *model.User, resultCode int, err error)
	}
)

var (
	localUserAuth IUserAuth
	localUserInfo IUserInfo
)

func UserAuth() IUserAuth {
	if localUserAuth == nil {
		panic("repository_implement localUserLogin not found for interface IUserAuth")
	}

	return localUserAuth
}

func InitUserAuth(i IUserAuth) {
	localUserAuth = i
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("repository_implement localUserInfo not found for interface IUserInfo")
	}

	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}
