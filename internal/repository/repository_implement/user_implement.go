package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"gorm.io/gorm"
	"time"
)

type rUser struct {
	db *gorm.DB
}

func NewUserRepositoryImplement(db *gorm.DB) *rUser {
	return &rUser{db: db}
}

func (r *rUser) CheckUserExistByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
	}

	return count > 0, nil
}

func (r *rUser) CreateUser(
	ctx context.Context,
	user *model.User,
) (*model.User, error) {
	res := r.db.WithContext(ctx).Create(user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *rUser) UpdateUser(
	ctx context.Context,
	userId uuid.UUID,
	updateData map[string]interface{},
) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *rUser) GetUser(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.User, error) {
	user := &model.User{}

	if res := r.db.WithContext(ctx).Model(user).Where(query, args...).First(user); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *rUser) GetManyUser(
	ctx context.Context,
	query *query_object.UserQueryObject,
) ([]*model.User, error) {
	var users []*model.User

	db := r.db.WithContext(ctx).Model(&model.User{})

	if query.Name != "" {
		db = db.Where("unaccent(family_name) ILIKE unaccent(?) OR unaccent(name) ILIKE unaccent(?)", "%"+query.Name+"%", "%"+query.Name+"%")
	}

	if query.Email != "" {
		db = db.Where("email = ?", query.Email)
	}

	if query.PhoneNumber != "" {
		db = db.Where("phonenumber = ?", query.PhoneNumber)
	}

	if !query.Birthday.IsZero() {
		birthday := query.Birthday.Truncate(24 * time.Hour)
		db = db.Where("birthday = ?", birthday)
	}

	if !query.CreatedAt.IsZero() {
		createAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createAt)
	}

	if query.SortBy != "" {
		switch query.SortBy {
		case "id":
			if query.IsDescending {
				db = db.Order("id DESC")
			} else {
				db = db.Order("id ASC")
			}
		case "name":
			combinedName := "unaccent(family_name || ' ' name)"
			if query.IsDescending {
				db = db.Order(combinedName + "DESC")
			} else {
				db = db.Order(combinedName + "ASC")
			}
		case "email":
			if query.IsDescending {
				db = db.Order("email DESC")
			} else {
				db = db.Order("email ASC")
			}
		case "phone_number":
			if query.IsDescending {
				db = db.Order("phone_number DESC")
			} else {
				db = db.Order("phone_number ASC")
			}
		case "birthday":
			if query.IsDescending {
				db = db.Order("birthday DESC")
			} else {
				db = db.Order("birthday ASC")
			}
		case "created_at":
			if query.IsDescending {
				db = db.Order("created_at DESC")
			} else {
				db = db.Order("created_at ASC")
			}
		}
	}

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
