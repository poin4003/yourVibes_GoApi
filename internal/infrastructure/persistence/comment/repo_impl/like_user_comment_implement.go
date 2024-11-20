package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/mapper"
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
	entity *entities.LikeUserComment,
) error {
	if err := r.db.WithContext(ctx).
		Create(mapper.ToLikeUserCommentModel(entity)).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rLikeUserComment) DeleteLikeUserComment(
	ctx context.Context,
	entity *entities.LikeUserComment,
) error {
	if err := r.db.WithContext(ctx).
		Delete(mapper.ToLikeUserCommentModel(entity)).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rLikeUserComment) GetLikeUserComment(
	ctx context.Context,
	query *query.GetCommentLikeQuery,
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

	err := db.Joins("JOIN like_user_comments ON like_user_comments.user_id = users.id").
		Where("like_user_comments.comment_id = ?", query.CommentId).
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

	var userEntities []*entities.User
	for _, user := range users {
		userEntities = append(userEntities, mapper.ToUserEntity(user))
	}

	return userEntities, pagingResponse, nil
}

func (r *rLikeUserComment) CheckUserLikeComment(
	ctx context.Context,
	entity *entities.LikeUserComment,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.LikeUserComment{}).
		Where("user_id=? AND comment_id=?", entity.UserId, entity.CommentId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
