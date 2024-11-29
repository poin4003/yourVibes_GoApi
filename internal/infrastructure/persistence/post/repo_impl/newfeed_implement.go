package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rNewFeed struct {
	db *gorm.DB
}

func NewNewFeedRepositoryImplement(db *gorm.DB) *rNewFeed {
	return &rNewFeed{db: db}
}

func (r *rNewFeed) CreateMany(
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

func (r *rNewFeed) DeleteOne(
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

func (r *rNewFeed) GetMany(
	ctx context.Context,
	query *query.GetNewFeedQuery,
) ([]*entities.Post, *response.PagingResponse, error) {
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
		Where("new_feeds.user_id = ?", query.UserId).
		Count(&total).Error

	if err != nil {
		return nil, nil, err
	}

	err = db.Model(&models.Post{}).
		Joins("JOIN new_feeds ON new_feeds.post_id = posts.id").
		Where("new_feeds.user_id = ?", query.UserId).
		Preload("User").
		Preload("Media").
		Order("posts.created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Error

	if err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postEntities []*entities.Post
	for _, post := range posts {
		postEntity := mapper.FromPostModel(post)
		postEntities = append(postEntities, postEntity)
	}

	return postEntities, pagingResponse, nil
}

func (r *rNewFeed) CreateManyWithRandomUser(
	ctx context.Context,
	numUsers int,
) error {
	query := `
		INSERT INTO new_feeds (user_id, post_id, view)
		SELECT u.id, a.post_id, 0
		FROM users u
		CROSS JOIN (
			SELECT advertises.id AS advertise_id, advertises.post_id
			FROM advertises
			JOIN bills ON bills.advertise_id = advertises.id
			WHERE bills.status = true
			AND advertises.start_date <= ?
			AND advertises.end_date >= ?
		) a
		WHERE NOT EXISTS (
			SELECT 1
			FROM new_feeds nf
			WHERE nf.user_id = u.id
			AND nf.post_id = a.post_id
		)
		ORDER BY RANDOM()
		LIMIT ?;
	`
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Exec(query, now, now, numUsers).Error; err != nil {
		return err
	}

	return nil
}

func (r *rNewFeed) DeleteExpiredAdvertiseFromNewFeeds(
	ctx context.Context,
) error {
	query := `
       	DELETE FROM new_feeds
		WHERE post_id IN (
			SELECT id 
			FROM posts
			WHERE is_advertisement = true 
			  AND EXISTS (
				  SELECT 1
				  FROM advertises
				  WHERE advertises.post_id = posts.id
					AND advertises.end_date < ?
			  )
		) 
    `

	if err := r.db.WithContext(ctx).
		Exec(query, time.Now()).
		Error; err != nil {
		return err
	}

	return nil
}
