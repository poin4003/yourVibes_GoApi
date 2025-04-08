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
) ([]*entities.UserWithSendFriendRequest, *response.PagingResponse, error) {
	var suggestions []*models.User
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

	// Total record of friend suggestion
	dbCount := r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("DISTINCT users.id").
		Joins("INNER JOIN friends f1 ON users.id = f1.friend_id").
		Joins("INNER JOIN friends f2 ON f1.user_id = f2.friend_id").
		Where("f2.user_id = ?", query.UserId).
		Where("users.id != ?", query.UserId).
		Where("users.id NOT IN (?)",
			r.db.Model(&models.Friend{}).
				Select("friend_id").
				Where("user_id = ?", query.UserId),
		)
	if err := dbCount.Count(&total).Error; err != nil {
		return nil, nil, response.NewServerFailedError("can not count friend suggestions")
	}

	// Get list of suggestion
	dbData := r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("DISTINCT users.id, users.family_name, users.name, users.avatar_url").
		Joins("INNER JOIN friends f1 ON users.id = f1.friend_id").
		Joins("INNER JOIN friends f2 ON f1.user_id = f2.friend_id").
		Where("f2.user_id = ?", query.UserId).
		Where("users.id != ?", query.UserId).
		Where("users.id NOT IN (?)",
			r.db.Model(&models.Friend{}).
				Select("friend_id").
				Where("user_id = ?", query.UserId),
		).
		Order("users.id").
		Limit(limit).
		Offset(offset)
	if err := dbData.Scan(&suggestions).Error; err != nil {
		return nil, nil, response.NewServerFailedError("can not get friend suggestions")
	}

	// Get UserId list from suggestion
	userIDs := make([]uuid.UUID, len(suggestions))
	for i, s := range suggestions {
		userIDs[i] = s.ID
	}

	// Check send friend request status
	var friendRequestResults []models.FriendRequest
	if len(userIDs) > 0 {
		if err := r.db.WithContext(ctx).
			Model(&models.FriendRequest{}).
			Select("friend_id").
			Where("user_id = ?", query.UserId).
			Where("friend_id IN (?)", userIDs).
			Find(&friendRequestResults).Error; err != nil {
			return nil, nil, response.NewServerFailedError("can not check friend request status")
		}
	}

	friendRequestStatus := make(map[uuid.UUID]bool)
	for _, fr := range friendRequestResults {
		friendRequestStatus[fr.FriendId] = true
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var userEntities []*entities.UserWithSendFriendRequest
	for _, sg := range suggestions {
		userEntity := mapper.FromUserModelWithSendFriendRequest(sg, friendRequestStatus[sg.ID])
		userEntities = append(userEntities, userEntity)
	}

	return userEntities, pagingResponse, nil
}
