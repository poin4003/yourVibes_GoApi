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
		VerifyEmail(ctx context.Context, email string) (err error)
		ChangePassword(ctx context.Context, command *command.ChangePasswordCommand) (err error)
		GetOtpForgotUserPassword(ctx context.Context, command *command.GetOtpForgotUserPasswordCommand) (err error)
		ForgotUserPassword(ctx context.Context, command *command.ForgotUserPasswordCommand) (err error)
		AuthGoogle(ctx context.Context, command *command.AuthGoogleCommand) (result *command.AuthGoogleCommandResult, err error)
		AppAuthGoogle(ctx context.Context, command *command.AuthAppGoogleCommand) (result *command.AuthGoogleCommandResult, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, query *query.GetOneUserQuery) (result *query.UserQueryResult, err error)
		GetUserStatusById(ctx context.Context, id uuid.UUID) (status bool, err error)
		GetManyUsers(ctx context.Context, query *query.GetManyUserQuery) (result *query.UserQueryListResult, err error)
		UpdateUser(ctx context.Context, command *command.UpdateUserCommand) (result *command.UpdateUserCommandResult, err error)
	}
	IUserFriend interface {
		SendAddFriendRequest(ctx context.Context, command *command.SendAddFriendRequestCommand) (err error)
		GetFriendRequests(ctx context.Context, query *query.FriendRequestQuery) (result *query.FriendRequestQueryResult, err error)
		AcceptFriendRequest(ctx context.Context, command *command.AcceptFriendRequestCommand) (err error)
		RemoveFriendRequest(ctx context.Context, command *command.RemoveFriendRequestCommand) (err error)
		UnFriend(ctx context.Context, command *command.UnFriendCommand) (err error)
		GetFriends(ctx context.Context, query *query.FriendQuery) (result *query.FriendQueryResult, err error)
		GetFriendSuggestion(ctx context.Context, query *query.FriendQuery) (result *query.FriendQueryResult, err error)
	}
)

var (
	localUserAuth   IUserAuth
	localUserInfo   IUserInfo
	localUserFriend IUserFriend
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

func InitUserFriend(i IUserFriend) {
	localUserFriend = i
}
