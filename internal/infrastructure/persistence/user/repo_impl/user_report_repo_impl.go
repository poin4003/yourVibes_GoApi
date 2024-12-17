package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rUserReport struct {
	db *gorm.DB
}

func NewUserReportRepositoryImplement(db *gorm.DB) *rUserReport {
	return &rUserReport{db: db}
}

func (r *rUserReport) GetById(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*entities.UserReport, error) {
	var userReportModel models.UserReport

	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Where("user_id = ? AND reported_user_id = ?", userId, reportedUserId).
		Preload("User").
		Preload("ReportedUser").
		Preload("Admin").
		First(&userReportModel).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromUserReportModel(&userReportModel), nil
}

func (r *rUserReport) CreateOne(
	ctx context.Context,
	entity *entities.UserReport,
) (*entities.UserReport, error) {
	userReportModel := mapper.ToUserReportModel(entity)

	if err := r.db.WithContext(ctx).
		Create(&userReportModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, entity.UserId, entity.ReportedUserId)
}

func (r *rUserReport) UpdateOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
	updateData *entities.UserReportUpdate,
) (*entities.UserReport, error) {
	updates := map[string]interface{}{}

	if updateData.AdminId != nil {
		updates["admin_id"] = *updateData.AdminId
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Where("user_id = ? AND reported_user_id = ?", userId, reportedUserId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, userId, reportedUserId)
}

func (r *rUserReport) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.UserReport{}).
		Where("user_id = ? AND reported_user_id = ?", userId, reportedUserId).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rUserReport) GetMany(
	ctx context.Context,
	query *query.GetManyUserReportQuery,
) ([]*entities.UserReport, *response.PagingResponse, error) {
	var userReportModels []*models.UserReport
	var total int64

	db := r.db.WithContext(ctx).Model(&models.UserReport{})

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
		case "reported_user_id":
			if query.IsDescending {
				db = db.Order("reported_user_id DESC")
			} else {
				db = db.Order("reported_user_id ASC")
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
		Find(&userReportModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var userReportEntities []*entities.UserReport
	for _, userReportModel := range userReportModels {
		userReportEntity := mapper.FromUserReportModel(userReportModel)
		userReportEntities = append(userReportEntities, userReportEntity)
	}

	return userReportEntities, pagingResponse, nil
}

func (r *rUserReport) CheckExist(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Where("user_id = ? AND reported_user_id = ?", userId, reportedUserId).
		Count(&count).Error; err != nil {
	}

	return count > 0, nil
}
