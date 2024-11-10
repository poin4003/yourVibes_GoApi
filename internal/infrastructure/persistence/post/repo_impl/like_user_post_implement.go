package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rLikeUserPost struct {
	db *gorm.DB
}

func NewLikeUserPostRepositoryImplement(db *gorm.DB) *rLikeUserPost {
	return &rLikeUserPost{db: db}
}

func (r *rLikeUserPost) CreateLikeUserPost(
	ctx context.Context,
	likeUserPost *models.LikeUserPost,
) error {
	res := r.db.WithContext(ctx).Create(likeUserPost)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rLikeUserPost) DeleteLikeUserPost(
	ctx context.Context,
	likeUserPost *models.LikeUserPost,
) error {
	res := r.db.WithContext(ctx).Delete(likeUserPost)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rLikeUserPost) GetLikeUserPost(
	ctx context.Context,
	postId uuid.UUID,
	query *query.PostLikeQueryObject,
) ([]*models.User, *response.PagingResponse, error) {
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

	err := db.Joins("JOIN like_user_posts ON like_user_posts.user_id = users.id").
		Where("like_user_posts.post_id = ?", postId).
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

func (r *rLikeUserPost) CheckUserLikePost(
	ctx context.Context,
	likeUserPost *models.LikeUserPost,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.LikeUserPost{}).
		Where("post_id = ? AND user_id =?", likeUserPost.PostId, likeUserPost.UserId).
		Count(&count).Error; err != nil {
	}
	return count > 0, nil
}
