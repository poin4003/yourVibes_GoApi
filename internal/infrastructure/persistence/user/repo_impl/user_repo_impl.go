package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
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
	id uuid.UUID,
) (*entities.User, error) {
	var userModel models.User
	if err := r.db.WithContext(ctx).
		Preload("Setting").
		First(&userModel, id).
		Error; err != nil {
		return nil, err
	}
	//fmt.Println(mapper.FromUserModel(&userModel))
	return mapper.FromUserModel(&userModel), nil
}

func (r *rUser) GetStatusById(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var userStatus bool
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("status").
		First(&userStatus, id).
		Error; err != nil {
		return false, err
	}
	return userStatus, nil
}

func (r *rUser) CreateOne(
	ctx context.Context,
	entity *entities.User,
) (*entities.User, error) {
	userModel := mapper.ToUserModel(entity)

	if err := r.db.WithContext(ctx).
		Create(userModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userModel.ID)
}

func (r *rUser) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.UserUpdate,
) (*entities.User, error) {
	updates := map[string]interface{}{}

	if updateData.FamilyName != nil {
		updates["family_name"] = *updateData.FamilyName
	}

	if updateData.Name != nil {
		updates["name"] = *updateData.Name
	}

	if updateData.PhoneNumber != nil {
		updates["phone_number"] = *updateData.PhoneNumber
	}

	if updateData.Birthday != nil {
		updates["birthday"] = *updateData.Birthday
	}

	if updateData.AvatarUrl != nil {
		updates["avatar_url"] = *updateData.AvatarUrl
	}

	if updateData.CapwallUrl != nil {
		updates["capwall_url"] = *updateData.CapwallUrl
	}

	if updateData.Privacy != nil {
		updates["privacy"] = *updateData.Privacy
	}

	if updateData.Biography != nil {
		updates["biography"] = *updateData.Biography
	}

	if updateData.AuthType != nil {
		updates["auth_type"] = *updateData.AuthType
	}

	if updateData.AuthGoogleId != nil {
		updates["auth_google_id"] = *updateData.AuthGoogleId
	}

	if updateData.PostCount != nil {
		updates["post_count"] = *updateData.PostCount
	}

	if updateData.FriendCount != nil {
		updates["friend_count"] = *updateData.FriendCount
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
	}

	if updateData.Password != nil {
		updates["password"] = *updateData.Password
	}

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rUser) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.User, error) {
	var userModel models.User

	if err := r.db.WithContext(ctx).
		Model(&userModel).
		Where(query, args...).
		Preload("Setting").
		First(&userModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userModel.ID)
}

func (r *rUser) GetMany(
	ctx context.Context,
	query *query.GetManyUserQuery,
) ([]*entities.User, *response.PagingResponse, error) {
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
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createdAt)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "id":
			sortColumn = "id"
		case "name":
			sortColumn = "unaccent(family_name || ' ' name)"
		case "email":
			sortColumn = "email"
		case "phone_number":
			sortColumn = "phone_number"
		case "birthday":
			sortColumn = "birthday"
		case "created_at":
			sortColumn = "created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
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
		Select("id, family_name, name, avatar_url").
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

	var userEntities []*entities.User
	for _, user := range userModels {
		userEntities = append(userEntities, mapper.FromUserModel(&user))
	}

	return userEntities, pagingResponse, nil
}

func (r *rUser) GetTotalUserCount(ctx context.Context) (int, error) {
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Count(&total).
		Error; err != nil {
		return 0, err
	}

	return int(total), nil
}
