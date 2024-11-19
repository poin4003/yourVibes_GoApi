package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	IUserRepository interface {
		CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
		GetById(ctx context.Context, id uuid.UUID) (*entities.User, error)
		CreateOne(ctx context.Context, entity *entities.User) (*entities.User, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.UserUpdate) (*entities.User, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.User, error)
		GetMany(ctx context.Context, query *query.GetManyUserQuery) ([]*entities.User, *response.PagingResponse, error)
	}
	ISettingRepository interface {
		GetById(ctx context.Context, id uint) (*entities.Setting, error)
		CreateOne(ctx context.Context, entity *entities.Setting) (*entities.Setting, error)
		UpdateOne(ctx context.Context, id uint, updateData *entities.SettingUpdate) (*entities.Setting, error)
		DeleteOne(ctx context.Context, id uint) error
		GetSetting(ctx context.Context, query interface{}, args ...interface{}) (*entities.Setting, error)
	}
	IFriendRequestRepository interface {
		CreateOne(ctx context.Context, entity *entities.FriendRequest) error
		DeleteOne(ctx context.Context, entity *entities.FriendRequest) error
		GetFriendRequests(ctx context.Context, query *query.FriendRequestQuery) ([]*entities.User, *response.PagingResponse, error)
		CheckFriendRequestExist(ctx context.Context, entity *entities.FriendRequest) (bool, error)
	}
	IFriendRepository interface {
		CreateOne(ctx context.Context, entity *entities.Friend) error
		DeleteOne(ctx context.Context, entity *entities.Friend) error
		GetFriends(ctx context.Context, query *query.FriendQuery) ([]*entities.User, *response.PagingResponse, error)
		GetFriendIds(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
		CheckFriendExist(ctx context.Context, entity *entities.Friend) (bool, error)
	}
)

var (
	localUser          IUserRepository
	localSetting       ISettingRepository
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

func InitFriendRequestRepository(i IFriendRequestRepository) {
	localFriendRequest = i
}

func InitFriendRepository(i IFriendRepository) {
	localFriend = i
}
