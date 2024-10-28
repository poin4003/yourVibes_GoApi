package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rFriend struct {
	db *gorm.DB
}

func NewFriendImplement(db *gorm.DB) *rMedia {
	return &rMedia{db: db}
}

func (r *rFriend) CreateFriend(
	ctx context.Context,
	friend *model.Friend,
) error {
	res := r.db.WithContext(ctx).Create(friend)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rFriend) DeleteFriend(
	ctx context.Context,
	friend *model.Friend,
) error {
	res := r.db.WithContext(ctx).Delete(friend)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriend) GetFriendRequest(
	ctx context.Context,
	userId uuid.UUID,
	query *query_object.FriendRequestQueryObject,
) ([]*model.User, *response.PagingResponse, error) {
	var users []*model.User
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

	db := r.db.WithContext(ctx).Model(&model.User{})

	err := db.Joins("JOIN friends ON friends.user_id = users.id").
		Where("friends.user_id = ?", userId).
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

	return users, pagingResponse, nil
}

func (r *rFriend) CheckFriendExist(
	ctx context.Context,
	friend *model.Friend,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&model.Friend{}).
		Where("friend_id = ? AND user_id = ?", friend.FriendId, friend.UserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
