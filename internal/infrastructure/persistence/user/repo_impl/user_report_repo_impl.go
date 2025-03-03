package repo_impl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/converter"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
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

func (r *rUserReport) UpdateMany(
	ctx context.Context,
	reportedUserId uuid.UUID,
	updateData *entities.UserReportUpdate,
) error {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return errors.New("no filed to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Where("reported_user_id = ?", reportedUserId).
		Updates(updates).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rUserReport) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND reported_user_id = ?", userId, reportedUserId).
		Delete(&models.UserReport{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rUserReport) DeleteByUserId(
	ctx context.Context,
	userId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Where("reported_user_id = ?", userId).
		Delete(&models.UserReport{}).
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

	db := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("ReportedUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})

	if query.Reason != "" {
		db = db.Where("user_reports.reason = ?", query.Reason)
	}

	if query.UserEmail != "" {
		db = db.Joins("JOIN users u ON user_reports.user_id = u.id").
			Where("u.email = ?", query.UserEmail)
	}

	if query.ReportedUserEmail != "" {
		db = db.Joins("JOIN users ru ON user_reports.reported_user_id = ru.id").
			Where("ru.email = ?", query.ReportedUserEmail)
	}

	if query.AdminEmail != "" {
		db = db.Joins("LEFT JOIN admins a ON user_reports.admin_id = a.id").
			Where("a.email = ?", query.AdminEmail)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("user_reports.created_at >= ?", query.FromDate)
	}
	if !query.ToDate.IsZero() {
		db = db.Where("user_reports.created_at <= ?", query.ToDate)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("user_reports.created_at = ?", createdAt)
	}

	if query.Status != nil {
		db = db.Where("user_reports.status = ?", *query.Status)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "user_id":
			sortColumn = "user_reports.user_id"
		case "reported_user_id":
			sortColumn = "user_reports.reported_user_id"
		case "admin_id":
			sortColumn = "user_reports.admin_id"
		case "reason":
			sortColumn = "user_reports.reason"
		case "created_at":
			sortColumn = "user_reports.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
		}
	}

	if err := db.Count(&total).Error; err != nil {
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
