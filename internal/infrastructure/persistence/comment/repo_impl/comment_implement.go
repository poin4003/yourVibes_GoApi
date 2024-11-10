package repo_impl

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
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

func (r *rComment) CreateComment(
	ctx context.Context,
	comment *models.Comment,
) (*models.Comment, error) {
	res := r.db.WithContext(ctx).Create(comment)

	if res.Error != nil {
		return nil, res.Error
	}

	return comment, nil
}

func (r *rComment) UpdateOneComment(
	ctx context.Context,
	commentId uuid.UUID,
	updateData map[string]interface{},
) (*models.Comment, error) {
	var comment models.Comment

	if err := r.db.WithContext(ctx).First(&comment, commentId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&comment).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *rComment) UpdateManyComment(
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

func (r *rComment) DeleteOneComment(
	ctx context.Context,
	commentId uuid.UUID,
) (*models.Comment, error) {
	comment := &models.Comment{}
	res := r.db.WithContext(ctx).First(comment, commentId)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Delete(comment)
	if res.Error != nil {
		return nil, res.Error
	}

	return comment, nil
}

func (r *rComment) DeleteManyComment(
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

func (r *rComment) GetOneComment(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*models.Comment, error) {
	comment := &models.Comment{}

	if res := r.db.WithContext(ctx).Model(comment).Where(query, args...).First(comment); res.Error != nil {
		return nil, res.Error
	}

	return comment, nil
}

func (r *rComment) GetManyComment(
	ctx context.Context,
	query *query.CommentQueryObject,
) ([]*models.Comment, *response.PagingResponse, error) {
	var comments []*models.Comment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Comment{})

	// 1. If query have ParentId
	if query.ParentId != "" {
		// 1.1. Find parent comment
		var parentComment models.Comment
		err := r.db.WithContext(ctx).Where("id = ?", query.ParentId).Find(&parentComment).Error
		if err != nil {
			return nil, nil, err
		}

		// 2.2. Find child comment by comment_left and comment_right of commentParent
		//db = db.Where("comment_left > ? AND comment_right <= ? ", parentComment.CommentLeft, parentComment.CommentRight)
		db = db.Where("parent_id = ?", parentComment.ID)
	} else if query.PostId != "" {
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

	return comments, pagingResponse, nil
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
