package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rFriendRequest struct {
	db *gorm.DB
}

func NewFriendRequestImplement(db *gorm.DB) *rFriendRequest {
	return &rFriendRequest{db: db}
}

func (r *rFriendRequest) CreateFriendRequest(
	ctx context.Context,
	friendRequest *model.FriendRequest,
) error {
	res := r.db.WithContext(ctx).Create(friendRequest)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rFriendRequest) DeleteFriendRequest(
	ctx context.Context,
	friendRequest *model.FriendRequest,
) error {
	res := r.db.WithContext(ctx).Delete(friendRequest)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rFriendRequest) GetFriendRequest(
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

	err := db.Joins("JOIN friend_requests ON friend_requests.user_id = users.id").
		Where("friend_requests.friend_id = ?", userId).
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

func (r *rFriendRequest) CheckFriendRequestExist(
	ctx context.Context,
	friendRequest *model.FriendRequest,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&model.FriendRequest{}).
		Where("friend_id = ? AND user_id = ?", friendRequest.FriendId, friendRequest.UserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
