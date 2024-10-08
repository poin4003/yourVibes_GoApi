package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, in *vo.LoginCredentials) (string, string, *model.User, error)
		Register(ctx context.Context, in *vo.RegisterCredentials) (int, error)
		VerifyEmail(ctx context.Context, email string) (int, error)
	}

	//IUserInfo interface {
	//	GetInfoByUserId(ctx context.Context) error
	//	GetAllUser(ctx context.Context) error
	//	FindOneUser(ctx context.Context) error
	//}
)

var (
	localUserAuth IUserAuth
	//localUserInfo  IUserInfo
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

//func UserInfo() IUserInfo {
//	if localUserInfo == nil {
//		panic("repository_implement localUserInfo not found for interface IUserInfo")
//	}
//
//	return localUserInfo
//}
//
//func InitUserInfo(i IUserInfo) {
//	localUserInfo = i
//}
