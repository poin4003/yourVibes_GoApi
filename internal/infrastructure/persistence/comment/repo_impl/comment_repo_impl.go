package repo_impl

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"strings"
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
		Where("status = true").
		Preload("User").
		First(&commentModel, id).
		Error; err != nil {
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
	updates := map[string]interface{}{}

	if updateData.Content != nil {
		updates["content"] = *updateData.Content
	}

	if updateData.LikeCount != nil {
		updates["like_count"] = *updateData.LikeCount
	}

	if updateData.RepCommentCount != nil {
		updates["rep_comment_count"] = *updateData.RepCommentCount
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
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

	for key, value := range condition {
		if strings.Contains(key, ">=") {
			db = db.Where(fmt.Sprintf("%s %s ?", key[:len(key)-2], ">="), value)
		} else if strings.Contains(key, ">") {
			db = db.Where(fmt.Sprintf("%s %s ?", key[:len(key)-1], ">"), value)
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

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
) error {
	db := r.db.WithContext(ctx).Model(&models.Comment{})

	for key, value := range condition {
		if strings.Contains(key, ">=") {
			db = db.Where(fmt.Sprintf("%s >= ?", key[:len(key)-3]), value)
		} else if strings.Contains(key, ">") {
			db = db.Where(fmt.Sprintf("%s > ?", key[:len(key)-2]), value)
		} else if strings.Contains(key, "<=") {
			db = db.Where(fmt.Sprintf("%s <= ?", key[:len(key)-3]), value)
		} else if strings.Contains(key, "<") {
			db = db.Where(fmt.Sprintf("%s < ?", key[:len(key)-2]), value)
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if err := db.Delete(condition).Error; err != nil {
		return err
	}

	return nil
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
		return nil, err
	}

	return r.GetById(ctx, comment.ID)
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

		// 2.2. Find child comment by comment_left and comment_right of commentParent
		//db = db.Where("comment_left > ? AND comment_right <= ? ", parentComment.CommentLeft, parentComment.CommentRight)
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

	if err := db.Offset(offset).Limit(limit).Order("comment_left ASC").Preload("User").Find(&comments).Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var commentEntities []*entities.Comment
	for _, commentModel := range comments {
		commentEntities = append(commentEntities, mapper.FromCommentModel(commentModel))
	}

	return commentEntities, pagingResponse, nil
}

func (r *rComment) GetMaxCommentRightByPostId(
	ctx context.Context,
	postId uuid.UUID,
) (int, error) {
	var maxRight *int
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("post_id = ?", postId).
		Select("MAX(comment_right)").
		Scan(&maxRight).Error

	if err != nil {
		return 0, err
	}

	if maxRight == nil {
		return 0, nil
	}

	return *maxRight, nil
}
