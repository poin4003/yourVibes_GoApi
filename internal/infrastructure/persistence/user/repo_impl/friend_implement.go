package repo_impl

import (
	"context"
	"github.com/google/uuid"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	user_mapper "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
	friendEntity *user_entity.Friend,
) error {
	friendModel := user_mapper.ToFriendModel(friendEntity)

	res := r.db.WithContext(ctx).Create(friendModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriend) DeleteOne(
	ctx context.Context,
	friendEntity *user_entity.Friend,
) error {
	friendModel := user_mapper.ToFriendModel(friendEntity)

	res := r.db.WithContext(ctx).Delete(friendModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriend) GetFriends(
	ctx context.Context,
	query *user_query.FriendQuery,
) ([]*user_entity.User, *response.PagingResponse, error) {
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
		Where("friends.user_id = ?", query.UserId).
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

	userEntities := user_mapper.FromUserModelList(users)

	return userEntities, pagingResponse, nil
}

func (r *rFriend) GetFriendIds(
	ctx context.Context,
	userId uuid.UUID,
) ([]uuid.UUID, error) {
	friendIds := []uuid.UUID{}

	err := r.db.WithContext(ctx).
		Model(&models.Friend{}).
		Where("user_id = ?", userId).
		Pluck("friend_id", &friendIds).Error

	if err != nil {
		return nil, err
	}

	return friendIds, nil
}

func (r *rFriend) CheckFriendExist(
	ctx context.Context,
	friendEntity *user_entity.Friend,
) (bool, error) {
	friend := user_mapper.ToFriendModel(friendEntity)
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.Friend{}).
		Where("friend_id = ? AND user_id = ?", friend.FriendId, friend.UserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
