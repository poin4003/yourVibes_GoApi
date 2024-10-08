package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

type (
	IUserRepository interface {
		CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
		CreateOne(ctx context.Context, user *model.User) (*model.User, error)
		UpdateOne(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}) (*model.User, error)
		GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error)
		GetUserByEmail(ctx context.Context, email string) (*model.User, error)
		GetAllUser(ctx context.Context) ([]*model.User, error)
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
