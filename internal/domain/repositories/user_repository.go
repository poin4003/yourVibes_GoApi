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
		GetStatusById(ctx context.Context, id uuid.UUID) (bool, error)
		CreateOne(ctx context.Context, entity *entities.User) (*entities.User, error)
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.UserUpdate) (*entities.User, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.User, error)
		GetMany(ctx context.Context, query *query.GetManyUserQuery) ([]*entities.User, *response.PagingResponse, error)
		GetTotalUserCount(ctx context.Context) (int, error)
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
	IUserReportRepository interface {
		GetById(ctx context.Context, userId uuid.UUID, reportedUserId uuid.UUID) (*entities.UserReport, error)
		CreateOne(ctx context.Context, entity *entities.UserReport) (*entities.UserReport, error)
		UpdateOne(ctx context.Context, userId uuid.UUID, reportedUserId uuid.UUID, updateData *entities.UserReportUpdate) (*entities.UserReport, error)
		UpdateMany(ctx context.Context, reportedUserId uuid.UUID, updateData *entities.UserReportUpdate) error
		DeleteOne(ctx context.Context, userId uuid.UUID, reportedUserId uuid.UUID) error
		DeleteByUserId(ctx context.Context, userId uuid.UUID) error
		GetMany(ctx context.Context, query *query.GetManyUserReportQuery) ([]*entities.UserReport, *response.PagingResponse, error)
		CheckExist(ctx context.Context, userId uuid.UUID, reportedUserId uuid.UUID) (bool, error)
	}
)

var (
	localUser          IUserRepository
	localSetting       ISettingRepository
	localFriendRequest IFriendRequestRepository
	localFriend        IFriendRepository
	localUserReport    IUserReportRepository
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

func UserReport() IUserReportRepository {
	if localUserReport == nil {
		panic("repository_implement localUserReport not found for interface IUserReport")
	}

	return localUserReport
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

func InitUserReportRepository(i IUserReportRepository) {
	localUserReport = i
}
