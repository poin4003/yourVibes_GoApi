package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rFriendRequest struct {
	db *gorm.DB
}

func NewFriendRequestImplement(db *gorm.DB) *rFriendRequest {
	return &rFriendRequest{db: db}
}

func (r *rFriendRequest) CreateOne(
	ctx context.Context,
	entity *entities.FriendRequest,
) error {
	friendRequestModel := mapper.ToFriendRequestModel(entity)

	res := r.db.WithContext(ctx).Create(friendRequestModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriendRequest) DeleteOne(
	ctx context.Context,
	entity *entities.FriendRequest,
) error {
	friendRequestModel := mapper.ToFriendRequestModel(entity)

	res := r.db.WithContext(ctx).Delete(friendRequestModel)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rFriendRequest) GetFriendRequests(
	ctx context.Context,
	query *query.FriendRequestQuery,
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

	err := db.Joins("JOIN friend_requests ON friend_requests.user_id = users.id").
		Where("friend_requests.friend_id = ?", query.UserId).
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

func (r *rFriendRequest) CheckFriendRequestExist(
	ctx context.Context,
	entity *entities.FriendRequest,
) (bool, error) {
	friendRequest := mapper.ToFriendRequestModel(entity)
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.FriendRequest{}).
		Where("friend_id = ? AND user_id = ?", friendRequest.FriendId, friendRequest.UserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
