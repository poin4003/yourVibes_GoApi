package repository_implement

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"gorm.io/gorm"
)

type rUser struct {
	db *gorm.DB
}

func NewUserRepositoryImplement(db *gorm.DB) *rUser {
	return &rUser{db: db}
}

func (r *rUser) CheckUserExistByEmail(ctx context.Context, email string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
	}

	return count > 0, nil
}

func (r *rUser) CreateOne(ctx context.Context, user *model.User) (*model.User, error) {
	res := r.db.WithContext(ctx).Create(user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *rUser) UpdateOne(ctx context.Context, userId uuid.UUID, updateData map[string]interface{}) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *rUser) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *rUser) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *rUser) GetAllUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
