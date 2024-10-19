package repository_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
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
	comment *model.Comment,
) (*model.Comment, error) {
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
) (*model.Comment, error) {
	var comment model.Comment

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
	db := r.db.WithContext(ctx).Model(&model.Comment{})

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

func (r *rComment) DeleteComment(
	ctx context.Context,
	commentId uuid.UUID,
) (*model.Comment, error) {
	comment := &model.Comment{}
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

func (r *rComment) GetComment(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Comment, error) {
	comment := &model.Comment{}

	if res := r.db.WithContext(ctx).Model(comment).Where(query, args...).First(comment); res.Error != nil {
		return nil, res.Error
	}

	return comment, nil
}

func (r *rComment) GetManyComment(
	ctx context.Context,
	query *query_object.CommentQueryObject,
) ([]*model.Comment, error) {
	var comments []*model.Comment

	db := r.db.WithContext(ctx).Model(&model.Comment{})

	// 1. If query have ParentId
	if query.ParentId != "" {
		// 1.1. Find parent comment
		var parentComment model.Comment
		err := r.db.WithContext(ctx).Where("id = ?", query.ParentId).Find(&parentComment).Error
		if err != nil {
			return nil, err
		}

		// 2.2. Find child comment by comment_left and comment_right of commentParent
		db = db.Where("comment_left > ? AND comment_right <= ?", parentComment.CommentLeft, parentComment.CommentRight)
	} else if query.PostId != "" {
		db = db.Where("post_id = ? AND parent_id IS NULL", query.PostId)
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

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Order("comment_left ASC").Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *rComment) GetMaxCommentRightByPostId(
	ctx context.Context,
	postId uuid.UUID,
) (int, error) {
	var maxRight *int
	err := r.db.WithContext(ctx).
		Model(&model.Comment{}).
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
