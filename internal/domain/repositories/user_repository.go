package repositories

import (
	"context"
	"github.com/google/uuid"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IUserRepository interface {
		CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
		GetById(ctx context.Context, userId uuid.UUID) (*user_entity.User, error)
		CreateOne(ctx context.Context, userEntity *user_entity.User) (*user_entity.User, error)
		UpdateOne(ctx context.Context, userId uuid.UUID, userUpdateEntity *user_entity.UserUpdate) (*user_entity.User, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*user_entity.User, error)
		GetMany(ctx context.Context, query *query.UserQueryObject) ([]*user_entity.User, *response.PagingResponse, error)
	}
	ISettingRepository interface {
		GetById(ctx context.Context, settingId uint) (*user_entity.Setting, error)
		CreateOne(ctx context.Context, settingEntity *user_entity.Setting) (*user_entity.Setting, error)
		UpdateOne(ctx context.Context, settingId uint, settingUpdateEntity *user_entity.SettingUpdate) (*user_entity.Setting, error)
		DeleteOne(ctx context.Context, settingId uint) error
		GetSetting(ctx context.Context, query interface{}, args ...interface{}) (*user_entity.Setting, error)
	}
	INotificationRepository interface {
		CreateOne(ctx context.Context, notificationEntity *user_entity.Notification) (*user_entity.Notification, error)
		CreateMany(ctx context.Context, notificationEntities []*user_entity.Notification) ([]*user_entity.Notification, error)
		UpdateOne(ctx context.Context, notificationId uint, updateData *user_entity.NotificationUpdate) (*user_entity.Notification, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, notificationId uint) (*user_entity.Notification, error)
		GetById(ctx context.Context, notificationId uint) (*user_entity.Notification, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*user_entity.Notification, error)
		GetMany(ctx context.Context, userId uuid.UUID, query *query.NotificationQueryObject) ([]*user_entity.Notification, *response.PagingResponse, error)
	}
	IFriendRequestRepository interface {
		CreateOne(ctx context.Context, friendRequestEntity *user_entity.FriendRequest) error
		DeleteOne(ctx context.Context, friendRequestEntity *user_entity.FriendRequest) error
		GetFriendRequests(ctx context.Context, query *user_query.FriendRequestQuery) ([]*user_entity.User, *response.PagingResponse, error)
		CheckFriendRequestExist(ctx context.Context, friendRequestEntity *user_entity.FriendRequest) (bool, error)
	}
	IFriendRepository interface {
		CreateOne(ctx context.Context, friendEntity *user_entity.Friend) error
		DeleteOne(ctx context.Context, friendEntity *user_entity.Friend) error
		GetFriends(ctx context.Context, query *user_query.FriendQuery) ([]*user_entity.User, *response.PagingResponse, error)
		GetFriendIds(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
		CheckFriendExist(ctx context.Context, friendEntity *user_entity.Friend) (bool, error)
	}
)

var (
	localUser          IUserRepository
	localSetting       ISettingRepository
	localNotification  INotificationRepository
	localFriendRequest IFriendRequestRepository
	localFriend        IFriendRepository
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
