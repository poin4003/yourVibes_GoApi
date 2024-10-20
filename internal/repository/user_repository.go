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
)

var (
	localUser    IUserRepository
	localSetting ISettingRepository
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

func InitUserRepository(i IUserRepository) {
	localUser = i
}

func InitSettingRepository(i ISettingRepository) {
	localSetting = i
}
