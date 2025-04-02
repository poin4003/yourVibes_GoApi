package repo_impl

import (
	"context"
	"errors"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"gorm.io/gorm"
)

type rFriend struct {
	db *gorm.DB
}

func NewFriendImplement(db *gorm.DB) *rFriend {
	return &rFriend{db: db}
}

func (r *rFriend) CreateOne(
	ctx context.Context,
	entity *entities.Friend,
) error {
	friendModel := mapper.ToFriendModel(entity)

	res := r.db.WithContext(ctx).Create(friendModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriend) DeleteOne(
	ctx context.Context,
	entity *entities.Friend,
) error {
	friendModel := mapper.ToFriendModel(entity)

	res := r.db.WithContext(ctx).Delete(friendModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriend) GetFriends(
	ctx context.Context,
	query *query.FriendQuery,
) ([]*entities.User, *response.PagingResponse, error) {
	var users []*models.User
	var total int64

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&models.User{})

	err := db.Joins("JOIN friends ON friends.user_id = users.id").
		Where("friends.friend_id = ?", query.UserId).
		Select("id, family_name, name, avatar_url").
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	userEntities := mapper.FromUserModelList(users)

	return userEntities, pagingResponse, nil
}

func (r *rFriend) GetFriendIds(
	ctx context.Context,
	userId uuid.UUID,
) ([]uuid.UUID, error) {
	friendIds := []uuid.UUID{}

	if err := r.db.WithContext(ctx).
		Model(&models.Friend{}).
		Where("user_id = ?", userId).
		Pluck("friend_id", &friendIds).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return friendIds, nil
}

func (r *rFriend) CheckFriendExist(
	ctx context.Context,
	entity *entities.Friend,
) (bool, error) {
	friend := mapper.ToFriendModel(entity)
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.Friend{}).
		Where("friend_id = ? AND user_id = ?", friend.FriendId, friend.UserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}

func (r *rFriend) GetFriendSuggestions(
	ctx context.Context,
	query *query.FriendQuery,
) ([]*entities.User, *response.PagingResponse, error) {
	var users []*models.User
	var total int64

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		INNER JOIN friends f1 ON u.id = f1.friend_id
		INNER JOIN friends f2 ON f1.user_id = f2.friend_id
		WHERE f2.user_id = ?
		AND u.id != ?
		AND u.id NOT IN (
			SELECT friend_id 
			FROM friends 
			WHERE user_id = ?
		)
	`
	if err := r.db.WithContext(ctx).Raw(countQuery, query.UserId, query.UserId, query.UserId).Scan(&total).Error; err != nil {
		return nil, nil, response.NewServerFailedError("can not count friend suggestions")
	}

	dataQuery := `
		SELECT DISTINCT u.id, u.family_name, u.name, u.avatar_url
		FROM users u
		INNER JOIN friends f1 ON u.id = f1.friend_id
		INNER JOIN friends f2 ON f1.user_id = f2.friend_id
		WHERE f2.user_id = ?
		AND u.id != ?
		AND u.id NOT IN (
			SELECT friend_id 
			FROM friends 
			WHERE user_id = ?
		)
		ORDER BY u.id 
		LIMIT ? OFFSET ?
	`
	if err := r.db.WithContext(ctx).Raw(dataQuery, query.UserId, query.UserId, query.UserId, limit, offset).Scan(&users).Error; err != nil {
		return nil, nil, response.NewServerFailedError("can not get friend suggestions")
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	userEntities := mapper.FromUserModelList(users)

	return userEntities, pagingResponse, nil
}
