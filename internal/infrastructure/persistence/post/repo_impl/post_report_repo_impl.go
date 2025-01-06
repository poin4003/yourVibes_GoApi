package repo_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
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
		Preload("ReportedPost.ParentPost.Media").
		Preload("ReportedPost.ParentPost.User").
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

func (r *rPostReport) UpdateMany(
	ctx context.Context,
	reportedPostId uuid.UUID,
	updateData *entities.PostReportUpdate,
) error {
	updates := map[string]interface{}{}

	if updateData.AdminId != nil {
		updates["admin_id"] = *updateData.AdminId
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Where("reported_post_id = ?", reportedPostId).
		Updates(updates).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rPostReport) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND reported_post_id = ?", userId, reportedPostId).
		Delete(&models.PostReport{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rPostReport) DeleteByPostId(
	ctx context.Context,
	postId uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Where("reported_post_id = ?", postId).
		Delete(&models.PostReport{}).
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

	db := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})

	if query.Reason != "" {
		db = db.Where("post_reports.reason = ?", query.Reason)
	}

	if query.UserEmail != "" {
		db = db.Joins("LEFT JOIN users ON users.id = post_reports.user_id").
			Where("users.email = ?", query.UserEmail)
	}

	if query.AdminEmail != "" {
		db = db.Joins("LEFT JOIN admins ON admins.id = post_reports.admin_id").
			Where("admins.email = ?", query.AdminEmail)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("post_reports.created_at >= ?", query.FromDate)
	}
	if !query.ToDate.IsZero() {
		db = db.Where("post_reports.created_at <= ?", query.ToDate)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createdAt)
	}

	if query.Status != nil {
		db = db.Where("post_reports.status = ?", *query.Status)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "user_id":
			sortColumn = "post_reports.user_id"
		case "reported_post_id":
			sortColumn = "post_reports.reported_post_id"
		case "admin_id":
			sortColumn = "post_reports.admin_id"
		case "reason":
			sortColumn = "post_reports.reason"
		case "created_at":
			sortColumn = "post_reports.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
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
