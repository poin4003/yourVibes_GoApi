package services

import (
	"context"
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IUserAuth interface {
		Login(ctx context.Context, loginCommand *user_command.LoginCommand) (result *user_command.LoginCommandResult, err error)
		Register(ctx context.Context, registerCommand *user_command.RegisterCommand) (result *user_command.RegisterCommandResult, err error)
		VerifyEmail(ctx context.Context, email string) (resultCode int, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context, getOneUserQuery *user_query.GetOneUserQuery) (result *user_query.UserQueryResult, err error)
		GetManyUsers(ctx context.Context, query *query.UserQueryObject) (result *user_query.UserQueryListResult, err error)
		UpdateUser(ctx context.Context, command *user_command.UpdateUserCommand) (result *user_command.UpdateUserCommandResult, err error)
	}
	IUserNotification interface {
		GetNotificationByUserId(ctx context.Context, userId uuid.UUID, query query.NotificationQueryObject) (notificationDtos []*response2.NotificationDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error)
		UpdateOneStatusNotification(ctx context.Context, notificationID uint) (resultCode int, httpStatusCode int, err error)
		UpdateManyStatusNotification(ctx context.Context, userId uuid.UUID) (resultCode int, httpStatusCode int, err error)
	}
	IUserFriend interface {
		SendAddFriendRequest(ctx context.Context, command *user_command.SendAddFriendRequestCommand) (result *user_command.SendAddFriendRequestCommandResult, err error)
		GetFriendRequests(ctx context.Context, query *user_query.FriendRequestQuery) (result *user_query.FriendRequestQueryResult, err error)
		AcceptFriendRequest(ctx context.Context, command *user_command.AcceptFriendRequestCommand) (result *user_command.AcceptFriendRequestCommandResult, err error)
		RemoveFriendRequest(ctx context.Context, command *user_command.RemoveFriendRequestCommand) (result *user_command.RemoveFriendRequestCommandResult, err error)
		UnFriend(ctx context.Context, command *user_command.UnFriendCommand) (result *user_command.UnFriendCommandResult, err error)
		GetFriends(ctx context.Context, query *user_query.FriendQuery) (result *user_query.FriendQueryResult, err error)
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
