package repo_impl

import (
	"context"
	"errors"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/mapper"
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
	commentFound := &models.Comment{}

	return mapper.FromCommentModel(commentFound), r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check if comment exists
		if err := tx.WithContext(ctx).
			Model(commentFound).
			First(commentFound, "id = ?", id).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		// 2. Delete related report
		result := tx.WithContext(ctx).
			Where("id IN (?)",
				tx.Model(&models.CommentReport{}).
					Select("report_id").
					Where("reported_comment_id = ?", id),
			).
			Delete(&models.Report{})

		if result.Error != nil {
			return response.NewServerFailedError(result.Error.Error())
		}

		// 3. Delete comment
		if err := tx.WithContext(ctx).
			Delete(&models.Comment{}, "id = ?", id).
			Error; err != nil {
			return err
		}
		return nil
	})
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
		var count int64
		if err := r.db.Where("id = ?", query.ParentId).
			Count(&count).
			Error; err != nil {
			return nil, nil, response.NewServerFailedError(err.Error())
		}

		if count == 0 {
			return nil, nil, response.NewDataNotFoundError("parent comment id does not exist")
		}

		db = db.Where("parent_id = ?", query.ParentId)
	} else if query.PostId != uuid.Nil {
		db = db.Where("post_id = ? AND parent_id IS NULL", query.PostId)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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

	if err = db.Offset(offset).Limit(limit).
		Order("created_at ASC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Find(&comments).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var commentEntities []*entities.Comment
	for _, comment := range comments {
		commentEntities = append(commentEntities, mapper.FromCommentModel(comment))
	}

	return commentEntities, pagingResponse, nil
}

func (r *rComment) DeleteCommentAndChildComment(
	ctx context.Context,
	commentId uuid.UUID,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Find comment
		comment := &models.Comment{}
		if err := tx.WithContext(ctx).
			Model(comment).
			Where("id = ?", commentId).
			Select("id, parent_id, post_id").
			First(comment).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Update -1 in parent comment
		if comment.ParentId != nil {
			if err := tx.WithContext(ctx).
				Model(&models.Comment{}).
				Where("id = ?", comment.ParentId).
				Update("rep_comment_count", gorm.Expr("rep_comment_count - ?", 1)).
				Error; err != nil {
				return response.NewServerFailedError(err.Error())
			}
		}

		// 3. Delete comment and child comment
		result := tx.WithContext(ctx).Exec(`
			WITH RECURSIVE cte AS (
				SELECT id FROM comments WHERE id = ?
				UNION ALL
				SELECT c.id FROM comments c INNER JOIN cte ON c.parent_id = cte.id
			)
			UPDATE comments 
			SET deleted_at = NOW() 
			WHERE id IN (SELECT id FROM cte) AND deleted_at IS NULL;
		`, commentId)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("no comments found to delete")
		}

		// 4. Delete related reports (CommentReport)
		if result.RowsAffected > 0 {
			commentReportResult := tx.WithContext(ctx).
				Where("id IN (?)",
					tx.Model(&models.CommentReport{}).
						Select("report_id").
						Where("reported_comment_id IN (?)",
							tx.Raw(`
                                WITH RECURSIVE cte AS (
                                    SELECT id FROM comments WHERE id = ?
                                    UNION ALL
                                    SELECT c.id FROM comments c INNER JOIN cte ON c.parent_id = cte.id
                                )
                                SELECT id FROM cte
                            `, commentId),
						),
				).
				Delete(&models.Report{})

			if commentReportResult.Error != nil {
				return response.NewServerFailedError(commentReportResult.Error.Error())
			}
		}

		// 5. Update commentCount in post table
		if err := tx.WithContext(ctx).
			Model(&models.Post{}).
			Where("id = ?", comment.PostId).
			Update("comment_count", gorm.Expr("comment_count - ?", result.RowsAffected)).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		return nil
	})
}
