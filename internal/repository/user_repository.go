package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

type (
	IUserRepository interface {
		CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
		CreateUser(ctx context.Context, user *model.User) (*model.User, error)
		UpdateUser(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}) (*model.User, error)
		GetUser(ctx context.Context, query interface{}, args ...interface{}) (*model.User, error)
		GetManyUser(ctx context.Context, limit, page int, query interface{}, args ...interface{}) ([]*model.User, int64, error)
	}
)

var (
	localUser IUserRepository
)

func User() IUserRepository {
	if localUser == nil {
		panic("repository_implement localUser not found for interface IUser")
	}

	return localUser
}

func InitUserRepository(i IUserRepository) {
	localUser = i
}
