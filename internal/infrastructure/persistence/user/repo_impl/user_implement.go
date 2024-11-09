package repo_impl

import (
	"context"
	"github.com/google/uuid"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	user_mapper "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).
		Error; err != nil {
	}

	return count > 0, nil
}

func (r *rUser) GetById(
	ctx context.Context,
	userId uuid.UUID,
) (*user_entity.User, error) {
	var userModel models.User
	if err := r.db.WithContext(ctx).
		First(&userModel, userId).
		Preload("Setting").
		Error; err != nil {
		return nil, err
	}
	return user_mapper.FromUserModel(&userModel), nil
}

func (r *rUser) CreateOne(
	ctx context.Context,
	userEntity *user_entity.User,
) (*user_entity.User, error) {
	userModel := user_mapper.ToUserModel(userEntity)

	if err := r.db.WithContext(ctx).
		Create(userModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userModel.ID)
}

func (r *rUser) UpdateOne(
	ctx context.Context,
	userId uuid.UUID,
	userUpdateEntity *user_entity.UserUpdate,
) (*user_entity.User, error) {
	updates := map[string]interface{}{}

	if userUpdateEntity.FamilyName != nil {
		updates["family_name"] = userUpdateEntity.FamilyName
	}

	if userUpdateEntity.Name != nil {
		updates["name"] = userUpdateEntity.Name
	}

	if userUpdateEntity.Email != nil {
		updates["email"] = userUpdateEntity.Email
	}

	if userUpdateEntity.Password != nil {
		updates["password"] = userUpdateEntity.Password
	}

	if userUpdateEntity.PhoneNumber != nil {
		updates["phone_number"] = userUpdateEntity.PhoneNumber
	}

	if userUpdateEntity.Birthday != nil {
		updates["birthday"] = userUpdateEntity.Birthday
	}

	if userUpdateEntity.AvatarUrl != nil {
		updates["avatar_url"] = userUpdateEntity.AvatarUrl
	}

	if userUpdateEntity.CapwallUrl != nil {
		updates["capwall_url"] = userUpdateEntity.CapwallUrl
	}

	if userUpdateEntity.Privacy != nil {
		updates["privacy"] = userUpdateEntity.Privacy
	}

	if userUpdateEntity.Biography != nil {
		updates["biography"] = userUpdateEntity.Biography
	}

	if userUpdateEntity.AuthType != nil {
		updates["auth_type"] = userUpdateEntity.AuthType
	}

	if userUpdateEntity.AuthGoogleId != nil {
		updates["auth_google_id"] = userUpdateEntity.AuthGoogleId
	}

	if userUpdateEntity.PostCount != nil {
		updates["post_count"] = userUpdateEntity.PostCount
	}

	if userUpdateEntity.FriendCount != nil {
		updates["friend_count"] = userUpdateEntity.FriendCount
	}

	if userUpdateEntity.Status != nil {
		updates["status"] = userUpdateEntity.Status
	}

	if userUpdateEntity.UpdatedAt != nil {
		updates["updated_at"] = userUpdateEntity.UpdatedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userId)
}

func (r *rUser) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*user_entity.User, error) {
	var userModel models.User

	if err := r.db.WithContext(ctx).
		Model(&userModel).
		Where(query, args...).
		Preload("Setting").
		First(userModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userModel.ID)
}

func (r *rUser) GetMany(
	ctx context.Context,
	query *query.UserQueryObject,
) ([]*user_entity.User, *response.PagingResponse, error) {
	var userModels []models.User
	var total int64

	db := r.db.WithContext(ctx).Model(&models.User{})

	if query.Name != "" {
		db = db.Where("unaccent(family_name || ' ' || name) ILIKE unaccent(?)", "%"+query.Name+"%")
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

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, err
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

	if err := db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&userModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	users := make([]*user_entity.User, len(userModels))
	for i, user := range userModels {
		users[i] = user_mapper.FromUserModel(&user)
	}

	return users, pagingResponse, nil
}
