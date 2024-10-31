package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IUserRepository interface {
		CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
		CreateUser(ctx context.Context, user *model.User) (*model.User, error)
		UpdateUser(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}) (*model.User, error)
		GetUser(ctx context.Context, query interface{}, args ...interface{}) (*model.User, error)
		GetManyUser(ctx context.Context, query *query_object.UserQueryObject) ([]*model.User, *response.PagingResponse, error)
	}
	ISettingRepository interface {
		CreateSetting(ctx context.Context, setting *model.Setting) (*model.Setting, error)
		UpdateSetting(ctx context.Context, settingId uint, updateData map[string]interface{}) (*model.Setting, error)
		DeleteSetting(ctx context.Context, settingId uint) error
		GetSetting(ctx context.Context, query interface{}, args ...interface{}) (*model.Setting, error)
	}
	INotificationRepository interface {
		CreateNotification(ctx context.Context, notification *model.Notification) (*model.Notification, error)
		CreateManyNotification(ctx context.Context, notifications []*model.Notification) ([]*model.Notification, error)
		UpdateOneNotification(ctx context.Context, notificationId uint, updateData map[string]interface{}) (*model.Notification, error)
		UpdateManyNotification(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteNotification(ctx context.Context, notificationId uint) (*model.Notification, error)
		GetOneNotification(ctx context.Context, query interface{}, args ...interface{}) (*model.Notification, error)
		GetManyNotification(ctx context.Context, userId uuid.UUID, query *query_object.NotificationQueryObject) ([]*model.Notification, *response.PagingResponse, error)
	}
	IFriendRequestRepository interface {
		CreateFriendRequest(ctx context.Context, friendRequest *model.FriendRequest) error
		DeleteFriendRequest(ctx context.Context, friendRequest *model.FriendRequest) error
		GetFriendRequest(ctx context.Context, userId uuid.UUID, query *query_object.FriendRequestQueryObject) ([]*model.User, *response.PagingResponse, error)
		CheckFriendRequestExist(ctx context.Context, friendRequest *model.FriendRequest) (bool, error)
	}
	IFriendRepository interface {
		CreateFriend(ctx context.Context, friend *model.Friend) error
		DeleteFriend(ctx context.Context, friend *model.Friend) error
		GetFriend(ctx context.Context, userId uuid.UUID, query *query_object.FriendQueryObject) ([]*model.User, *response.PagingResponse, error)
		GetFriendIds(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
		CheckFriendExist(ctx context.Context, friend *model.Friend) (bool, error)
	}
	INewFeedRepository interface {
		CreateManyNewFeed(ctx context.Context, postId uuid.UUID, friendIds []uuid.UUID) error
		DeleteNewFeed(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
		GetManyNewFeed(ctx context.Context, userId uuid.UUID, query *query_object.NewFeedQueryObject) ([]*model.Post, *response.PagingResponse, error)
	}
)

var (
	localUser          IUserRepository
	localSetting       ISettingRepository
	localNotification  INotificationRepository
	localFriendRequest IFriendRequestRepository
	localFriend        IFriendRepository
	localNewFeed       INewFeedRepository
)

func User() IUserRepository {
	if localUser == nil {
		panic("repository_implement localUser not found for interface IUser")
	}

	return localUser
}

func Setting() ISettingRepository {
	if localSetting == nil {
		panic("repository_implement localSetting not found for interface ISetting")
	}

	return localSetting
}

func Notification() INotificationRepository {
	if localNotification == nil {
		panic("repository_implement localNotification not found for interface INotification")
	}

	return localNotification
}

func FriendRequest() IFriendRequestRepository {
	if localFriendRequest == nil {
		panic("repository_implement localFriendRequest not found for interface IFriendRequest")
	}

	return localFriendRequest
}

func NewFeed() INewFeedRepository {
	if localNewFeed == nil {
		panic("repository_implement localNewFeed not found for interface INewFeed")
	}

	return localNewFeed
}

func Friend() IFriendRepository {
	if localFriend == nil {
		panic("repository_implement localFriendRequest not found for interface IFriend")
	}

	return localFriend
}

func InitUserRepository(i IUserRepository) {
	localUser = i
}

func InitSettingRepository(i ISettingRepository) {
	localSetting = i
}

func InitNotificationRepository(i INotificationRepository) {
	localNotification = i
}

func InitFriendRequestRepository(i IFriendRequestRepository) {
	localFriendRequest = i
}

func InitFriendRepository(i IFriendRepository) {
	localFriend = i
}

func InitNewFeedRepository(i INewFeedRepository) {
	localNewFeed = i
}
