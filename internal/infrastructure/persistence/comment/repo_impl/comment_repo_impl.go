package repo_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/converter"
	"gorm.io/gorm"
)

type rComment struct {
	db *gorm.DB
}

func NewCommentRepositoryImplement(db *gorm.DB) *rComment {
	return &rComment{db: db}
}

func (r *rComment) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Comment, error) {
	var commentModel models.Comment
	if err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		First(&commentModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.FromCommentModel(&commentModel), nil
}

func (r *rComment) CreateOne(
	ctx context.Context,
	entity *entities.Comment,
) (*entities.Comment, error) {
	commentModel := mapper.ToCommentModel(entity)

	if err := r.db.WithContext(ctx).
		Create(commentModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, entity.ID)
}

func (r *rComment) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.CommentUpdate,
) (*entities.Comment, error) {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, nil
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rComment) UpdateMany(
	ctx context.Context,
	condition map[string]interface{},
	updateData map[string]interface{},
) error {
	db := r.db.WithContext(ctx).Model(&models.Comment{})

	if err := db.Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *rComment) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Comment, error) {
	comment, err := r.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := r.db.WithContext(ctx).Delete(comment)
	if res.Error != nil {
		return nil, res.Error
	}

	return comment, nil
}

func (r *rComment) DeleteMany(
	ctx context.Context,
	condition map[string]interface{},
) (int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Comment{})

	result := db.Delete(condition)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *rComment) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Comment, error) {
	comment := &models.Comment{}

	if err := r.db.WithContext(ctx).
		Model(comment).
		Where(query, args...).
		Where("status = true").
		First(comment).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.FromCommentModel(comment), nil
}

func (r *rComment) GetMany(
	ctx context.Context,
	query *query.GetManyCommentQuery,
) ([]*entities.Comment, *response.PagingResponse, error) {
	var comments []*models.Comment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Comment{})

	db = db.Where("status = true")

	// 1. If query have ParentId
	if query.ParentId != uuid.Nil {
		// 1.1. Find parent comment
		var parentComment models.Comment
		if err := r.db.Where("id = ?", query.ParentId).
			Find(&parentComment).
			Error; err != nil {
			return nil, nil, err
		}

		db = db.Where("parent_id = ?", parentComment.ID)
	} else if query.PostId != uuid.Nil {
		db = db.Where("post_id = ? AND parent_id IS NULL", query.PostId)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := db.Offset(offset).Limit(limit).
		Order("comment_left ASC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Find(&comments).
		Error; err != nil {
		return nil, nil, err
	}

	var commentIds []uuid.UUID
	for _, comment := range comments {
		commentIds = append(commentIds, comment.ID)
	}

	var likeCommentIds []uuid.UUID
	if err := r.db.Model(&models.LikeUserComment{}).
		Select("comment_id").
		Where("user_id = ? AND comment_id IN ?", query.AuthenticatedUserId, commentIds).
		Find(&likeCommentIds).
		Error; err != nil {
		return nil, nil, err
	}

	likedMap := make(map[uuid.UUID]bool)
	for _, id := range likeCommentIds {
		likedMap[id] = true
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var commentEntities []*entities.Comment
	for _, comment := range comments {
		commentEntities = append(commentEntities, mapper.FromCommentModelWithLiked(comment, likedMap[comment.ID]))
	}

	return commentEntities, pagingResponse, nil
}

func (r *rComment) DeleteCommentAndChildComment(
	ctx context.Context,
	commentId uuid.UUID,
) (int64, error) {
	var deleteCount int64

	tx := r.db.WithContext(ctx).Exec(`
		WITH RECURSIVE cte AS (
			SELECT id FROM comments WHERE id = ?
			UNION ALL
			SELECT c.id FROM comments c INNER JOIN cte ON c.parent_id = cte.id
		)
		UPDATE comments 
		SET deleted_at = NOW() 
		WHERE id IN (SELECT id FROM cte) AND deleted_at IS NULL;
	`, commentId)

	if tx.Error != nil {
		return 0, tx.Error
	}

	deleteCount = tx.RowsAffected
	if deleteCount == 0 {
		return 0, errors.New("no comments found to delete")
	}

	return deleteCount, nil
}
