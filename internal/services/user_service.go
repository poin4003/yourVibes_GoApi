package services

//
//import (
//	"context"
//	"github.com/poin4003/yourVibes_GoApi/internal/model"
//)
//
//type (
//	IUserLogin interface {
//		Login(ctx context.Context, in *model.LoginInput) error
//		Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error)
//		VerifyOTP(ctx context.Context, in *model.VerifyOtpInput) error
//	}
//
//	IUserInfo interface {
//		GetInfoByUserId(ctx context.Context) error
//		GetAllUser(ctx context.Context) error
//		FindOneUser(ctx context.Context) error
//	}
//)
//
//var (
//	localUserLogin IUserLogin
//	localUserInfo  IUserInfo
//)
//
//func UserLogin() IUserLogin {
//	if localUserLogin == nil {
//		panic("implement localUserLogin not found for interface IUserLogin")
//	}
//
//	return localUserLogin
//}
//
//func InitUserLogin(i IUserLogin) {
//	localUserLogin = i
//}
//
//func UserInfo() IUserInfo {
//	if localUserInfo == nil {
//		panic("implement localUserInfo not found for interface IUserInfo")
//	}
//
//	return localUserInfo
//}
//
//func InitUserInfo(i IUserInfo) {
//	localUserInfo = i
//}
