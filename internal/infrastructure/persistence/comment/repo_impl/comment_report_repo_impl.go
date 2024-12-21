package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rCommentReport struct {
	db *gorm.DB
}

func NewCommentReportRepositoryImplement(db *gorm.DB) *rCommentReport {
	return &rCommentReport{db: db}
}

func (r *rCommentReport) GetById(
	ctx context.Context,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*entities.CommentReport, error) {
	var commentReportModel models.CommentReport

	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Where("user_id = ? AND reported_comment_id = ?", userId, reportedCommentId).
		Preload("User").
		Preload("ReportedComment.Post.User").
		Preload("ReportedComment.User").
		Preload("Admin").
		First(&commentReportModel).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromCommentReportModel(&commentReportModel), nil
}

func (r *rCommentReport) CreateOne(
	ctx context.Context,
	entity *entities.CommentReport,
) (*entities.CommentReport, error) {
	commentReportModel := mapper.ToCommentReportModel(entity)

	if err := r.db.WithContext(ctx).
		Create(&commentReportModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, entity.UserId, entity.ReportedCommentId)
}

func (r *rCommentReport) UpdateOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
	updateData *entities.CommentReportUpdate,
) (*entities.CommentReport, error) {
	updates := map[string]interface{}{}

	if updateData.AdminId != nil {
		updates["admin_id"] = *updateData.AdminId
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Where("user_id = ? AND reported_comment_id = ?", userId, reportedCommentId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userId, reportedCommentId)
}

func (r *rCommentReport) UpdateMany(
	ctx context.Context,
	reportedCommentId uuid.UUID,
	updateData *entities.CommentReportUpdate,
) error {
	updates := map[string]interface{}{}

	if updateData.AdminId != nil {
		updates["admin_id"] = *updateData.AdminId
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Where("reported_comment_id = ?", reportedCommentId).
		Updates(updates).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rCommentReport) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND reported_comment_id = ?", userId, reportedCommentId).
		Delete(&models.CommentReport{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rCommentReport) GetMany(
	ctx context.Context,
	query *query.GetManyCommentReportQuery,
) ([]*entities.CommentReport, *response.PagingResponse, error) {
	var commentReportModels []*models.CommentReport
	var total int64

	db := r.db.WithContext(ctx).Model(&models.CommentReport{})

	if query.Reason != "" {
		db = db.Where("reason = ?", query.Reason)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createdAt)
	}

	if query.Status != nil {
		if *query.Status {
			db = db.Where("status = ?", true)
		} else {
			db = db.Where("status = ?", false)
		}
	}

	if query.SortBy != "" {
		switch query.SortBy {
		case "user_id":
			if query.IsDescending {
				db = db.Order("user_id DESC")
			} else {
				db = db.Order("user_id ASC")
			}
		case "reported_comment_id":
			if query.IsDescending {
				db = db.Order("reported_comment_id DESC")
			} else {
				db = db.Order("reported_comment_id ASC")
			}
		case "admin_id":
			if query.IsDescending {
				db = db.Order("admin_id DESC")
			} else {
				db = db.Order("admin_id ASC")
			}
		case "reason":
			if query.IsDescending {
				db = db.Order("reason DESC")
			} else {
				db = db.Order("reason ASC")
			}
		case "created_at":
			if query.IsDescending {
				db = db.Order("created_at DESC")
			} else {
				db = db.Order("created_at ASC")
			}
		}
	}

	if err := db.Count(&total).
		Error; err != nil {
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

	if err := db.Offset(offset).
		Limit(limit).
		Find(&commentReportModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var commentReportEntities []*entities.CommentReport
	for _, commentReportModel := range commentReportModels {
		commentReportEntity := mapper.FromCommentReportModel(commentReportModel)
		commentReportEntities = append(commentReportEntities, commentReportEntity)
	}

	return commentReportEntities, pagingResponse, nil
}

func (r *rCommentReport) CheckExist(
	ctx context.Context,
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Where("user_id = ? AND reported_comment_id = ?", userId, reportedCommentId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
