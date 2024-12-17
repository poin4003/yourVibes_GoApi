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

type rPostReport struct {
	db *gorm.DB
}

func NewPostReportRepositoryImplement(db *gorm.DB) *rPostReport {
	return &rPostReport{db: db}
}

func (r *rPostReport) GetById(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*entities.PostReport, error) {
	var postReportModel models.PostReport

	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Where("user_id = ? AND reported_post_id = ?", userId, reportedPostId).
		Preload("User").
		Preload("ReportedPost.User").
		Preload("Admin").
		First(&postReportModel).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromPostReportModel(&postReportModel), nil
}

func (r *rPostReport) CreateOne(
	ctx context.Context,
	entity *entities.PostReport,
) (*entities.PostReport, error) {
	postReportModel := mapper.ToPostReportModel(entity)

	if err := r.db.WithContext(ctx).
		Create(&postReportModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, entity.UserId, entity.ReportedPostId)
}

func (r *rPostReport) UpdateOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
	updateData *entities.PostReportUpdate,
) (*entities.PostReport, error) {
	updates := map[string]interface{}{}

	if updateData.AdminId != nil {
		updates["admin_id"] = *updateData.AdminId
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Where("user_id = ? AND reported_post_id = ?", userId, reportedPostId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userId, reportedPostId)
}

func (r *rPostReport) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.PostReport{}).
		Where("user_id = ? AND reported_post_id = ?", userId, reportedPostId).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rPostReport) GetMany(
	ctx context.Context,
	query *query.GetManyPostReportQuery,
) ([]*entities.PostReport, *response.PagingResponse, error) {
	var postReportModels []*models.PostReport
	var total int64

	db := r.db.WithContext(ctx).Model(&models.PostReport{})

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
		case "reported_post_id":
			if query.IsDescending {
				db = db.Order("reported_post_id DESC")
			} else {
				db = db.Order("reported_post_id ASC")
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
		Find(&postReportModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postReportEntities []*entities.PostReport
	for _, postReportModel := range postReportModels {
		postReportEntity := mapper.FromPostReportModel(postReportModel)
		postReportEntities = append(postReportEntities, postReportEntity)
	}

	return postReportEntities, pagingResponse, nil
}

func (r *rPostReport) CheckExist(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Where("user_id = ? AND reported_post_id = ?", userId, reportedPostId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
