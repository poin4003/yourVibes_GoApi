package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rLikeUserComment struct {
	db *gorm.DB
}

func NewLikeUserCommentRepositoryImplement(db *gorm.DB) *rLikeUserComment {
	return &rLikeUserComment{db: db}
}

func (r *rLikeUserComment) CreateLikeUserComment(
	ctx context.Context,
	likeUserComment *models.LikeUserComment,
) error {
	res := r.db.WithContext(ctx).Create(likeUserComment)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rLikeUserComment) DeleteLikeUserComment(
	ctx context.Context,
	likeUserComment *models.LikeUserComment,
) error {
	res := r.db.WithContext(ctx).Delete(likeUserComment)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rLikeUserComment) GetLikeUserComment(
	ctx context.Context,
	commentId uuid.UUID,
	query *query.CommentLikeQueryObject,
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

	err := db.Joins("JOIN like_user_comments ON like_user_comments.user_id = users.id").
		Where("like_user_comments.comment_id = ?", commentId).
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

func (r *rLikeUserComment) CheckUserLikeComment(
	ctx context.Context,
	likeUserComment *models.LikeUserComment,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.LikeUserComment{}).
		Where("user_id=? AND comment_id=?", likeUserComment.UserId, likeUserComment.CommentId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}