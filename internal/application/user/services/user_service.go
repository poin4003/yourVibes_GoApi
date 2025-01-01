package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, loginCommand *command.LoginCommand) (result *command.LoginCommandResult, err error)
		Register(ctx context.Context, registerCommand *command.RegisterCommand) (result *command.RegisterCommandResult, err error)
		VerifyEmail(ctx context.Context, email string) (resultCode int, err error)
		AuthGoogle(ctx context.Context, command *command.AuthGoogleCommand) (result *command.AuthGoogleCommandResult, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, query *query.GetOneUserQuery) (result *query.UserQueryResult, err error)
		GetUserStatusById(ctx context.Context, id uuid.UUID) (status bool, err error)
		GetManyUsers(ctx context.Context, query *query.GetManyUserQuery) (result *query.UserQueryListResult, err error)
		UpdateUser(ctx context.Context, command *command.UpdateUserCommand) (result *command.UpdateUserCommandResult, err error)
	}
	IUserNotification interface {
		GetNotificationByUserId(ctx context.Context, query *query.GetManyNotificationQuery) (result *query.GetManyNotificationQueryResult, err error)
		UpdateOneStatusNotification(ctx context.Context, command *command.UpdateOneStatusNotificationCommand) (result *command.UpdateOneStatusNotificationCommandResult, err error)
		UpdateManyStatusNotification(ctx context.Context, command *command.UpdateManyStatusNotificationCommand) (result *command.UpdateManyStatusNotificationCommandResult, err error)
	}
	IUserFriend interface {
		SendAddFriendRequest(ctx context.Context, command *command.SendAddFriendRequestCommand) (result *command.SendAddFriendRequestCommandResult, err error)
		GetFriendRequests(ctx context.Context, query *query.FriendRequestQuery) (result *query.FriendRequestQueryResult, err error)
		AcceptFriendRequest(ctx context.Context, command *command.AcceptFriendRequestCommand) (result *command.AcceptFriendRequestCommandResult, err error)
		RemoveFriendRequest(ctx context.Context, command *command.RemoveFriendRequestCommand) (result *command.RemoveFriendRequestCommandResult, err error)
		UnFriend(ctx context.Context, command *command.UnFriendCommand) (result *command.UnFriendCommandResult, err error)
		GetFriends(ctx context.Context, query *query.FriendQuery) (result *query.FriendQueryResult, err error)
	}
	IUserReport interface {
		CreateUserReport(ctx context.Context, command *command.CreateReportUserCommand) (result *command.CreateReportUserCommandResult, err error)
		HandleUserReport(ctx context.Context, command *command.HandleUserReportCommand) (result *command.HandleUserReportCommandResult, err error)
		DeleteUserReport(ctx context.Context, command *command.DeleteUserReportCommand) (result *command.DeleteUserReportCommandResult, err error)
		ActivateUserAccount(ctx context.Context, command *command.ActivateUserAccountCommand) (result *command.ActivateUserAccountCommandResult, err error)
		GetDetailUserReport(ctx context.Context, query *query.GetOneUserReportQuery) (result *query.UserReportQueryResult, err error)
		GetManyUserReport(ctx context.Context, query *query.GetManyUserReportQuery) (result *query.UserReportQueryListResult, err error)
	}
)

var (
	localUserAuth         IUserAuth
	localUserInfo         IUserInfo
	localUserNotification IUserNotification
	localUserFriend       IUserFriend
	localUserReport       IUserReport
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

func UserFriend() IUserFriend {
	if localUserFriend == nil {
		panic("repository_implement localUserFriend not found for interface IUserFriend")
	}

	return localUserFriend
}

func UserReport() IUserReport {
	if localUserReport == nil {
		panic("repository_implement localUserReport not found for interface IUserReport")
	}

	return localUserReport
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

func InitUserFriend(i IUserFriend) {
	localUserFriend = i
}

func InitUserReport(i IUserReport) {
	localUserReport = i
}
