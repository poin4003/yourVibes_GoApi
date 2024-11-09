package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rNewFeed struct {
	db *gorm.DB
}

func NewNewFeedRepositoryImplement(db *gorm.DB) *rNewFeed {
	return &rNewFeed{db: db}
}

func (r *rNewFeed) CreateManyNewFeed(
	ctx context.Context,
	postId uuid.UUID,
	friendIds []uuid.UUID,
) error {
	// 1. Create single post for single friend
	var newFeeds []models.NewFeed
	for _, friendId := range friendIds {
		newFeeds = append(newFeeds, models.NewFeed{
			UserId: friendId,
			PostId: postId,
			View:   0,
		})
	}

	// 2. Create new feed in db
	if err := r.db.WithContext(ctx).Create(&newFeeds).Error; err != nil {
		return err
	}
	return nil
}

func (r *rNewFeed) DeleteNewFeed(
	ctx context.Context,
	userId uuid.UUID,
	postId uuid.UUID,
) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userId, postId).
		Delete(&models.NewFeed{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rNewFeed) GetManyNewFeed(
	ctx context.Context,
	userId uuid.UUID,
	query *query.NewFeedQueryObject,
) ([]*models.Post, *response.PagingResponse, error) {
	var posts []*models.Post
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

	db := r.db.WithContext(ctx)

	err := db.Model(&models.Post{}).
		Joins("JOIN new_feeds ON new_feeds.post_id = posts.id").
		Where("new_feeds.user_id = ?", userId).
		Count(&total).Error

	if err != nil {
		return nil, nil, err
	}

	err = db.Model(&models.Post{}).
		Joins("JOIN new_feeds ON new_feeds.post_id = posts.id").
		Where("new_feeds.user_id = ?", userId).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return posts, pagingResponse, nil
}
