package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, loginCommand *command.LoginCommand) (result *command.LoginCommandResult, err error)
		Register(ctx context.Context, registerCommand *command.RegisterCommand) (result *command.RegisterCommandResult, err error)
		VerifyEmail(ctx context.Context, email string) (resultCode int, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, query *query.GetOneUserQuery) (result *query.UserQueryResult, err error)
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
)

var (
	localUserAuth         IUserAuth
	localUserInfo         IUserInfo
	localUserNotification IUserNotification
	localUserFriend       IUserFriend
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
