package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rNotification struct {
	db *gorm.DB
}

func NewNotificationRepositoryImplement(db *gorm.DB) *rNotification {
	return &rNotification{db: db}
}

func (r *rNotification) GetById(
	ctx context.Context,
	id uint,
) (*entities.Notification, error) {
	var notificationModel models.Notification
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&notificationModel, id).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromNotificationModel(&notificationModel), nil
}

func (r *rNotification) CreateOne(
	ctx context.Context,
	entity *entities.Notification,
) (*entities.Notification, error) {
	notificationModel := mapper.ToNotificationModel(entity)

	res := r.db.WithContext(ctx).Create(notificationModel)

	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, notificationModel.ID)
}

func (r *rNotification) CreateMany(
	ctx context.Context,
	notificationEntities []*entities.Notification,
) ([]*entities.Notification, error) {
	var notificationModels []*models.Notification
	for i, notification := range notificationEntities {
		notificationModels[i] = mapper.ToNotificationModel(notification)
	}

	err := r.db.WithContext(ctx).Create(&notificationModels).Error
	if err != nil {
		return nil, err
	}

	notificationEntityList := make([]*entities.Notification, len(notificationModels))
	for i, notificationEntity := range notificationModels {
		notificationEntityList[i] = mapper.FromNotificationModel(notificationEntity)
	}

	return notificationEntityList, nil
}

func (r *rNotification) UpdateOne(
	ctx context.Context,
	notificationId uint,
	updateData *entities.NotificationUpdate,
) (*entities.Notification, error) {
	updates := map[string]interface{}{}

	if updateData.From != nil {
		updates["from"] = updateData.From
	}

	if updateData.FromUrl != nil {
		updates["from_url"] = updateData.FromUrl
	}

	if updateData.NotificationType != nil {
		updates["notification_type"] = updateData.NotificationType
	}

	if updateData.ContentId != nil {
		updates["content_id"] = updateData.ContentId
	}

	if updateData.Content != nil {
		updates["content"] = updateData.Content
	}

	if updateData.Status != nil {
		updates["status"] = updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = updateData.UpdatedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", notificationId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, notificationId)
}

func (r *rNotification) UpdateMany(
	ctx context.Context,
	condition map[string]interface{},
	updateData map[string]interface{},
) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where(condition).
		Updates(updateData).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rNotification) DeleteOne(
	ctx context.Context,
	id uint,
) (*entities.Notification, error) {
	res := r.db.WithContext(ctx).
		Delete(&models.Notification{}, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, id)
}

func (r *rNotification) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Notification, error) {
	var notificationModel models.Notification

	if res := r.db.WithContext(ctx).
		Model(&notificationModel).
		Where(query, args...).
		Preload("User").
		First(&notificationModel); res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, notificationModel.ID)
}

func (r *rNotification) GetMany(
	ctx context.Context,
	query *query.GetManyNotificationQuery,
) ([]*entities.Notification, *response.PagingResponse, error) {
	var notificationModels []*models.Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Notification{})

	if query.From != "" {
		db = db.Where("LOWER(from) LIKE LOWER(?)", "%"+query.From+"%")
	}

	if query.NotificationType != "" {
		db = db.Where("LOWER(notification_type) LIKE LOWER(?)", "%"+query.NotificationType+"%")
	}

	if !query.CreatedAt.IsZero() {
		createAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createAt)
	}

	if query.SortBy != "" {
		switch query.SortBy {
		case "id":
			if query.IsDescending {
				db = db.Order("id DESC")
			} else {
				db = db.Order("id ASC")
			}
		case "from":
			if query.IsDescending {
				db = db.Order("from DESC")
			} else {
				db = db.Order("from ASC")
			}
		case "notification_type":
			if query.IsDescending {
				db = db.Order("notification_type DESC")
			} else {
				db = db.Order("notification_type ASC")
			}
		case "created_at":
			if query.IsDescending {
				db = db.Order("created_at DESC")
			} else {
				db = db.Order("created_at ASC")
			}
		}
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

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).
		Where("user_id=?", query.UserId).
		Preload("User").
		Find(&notificationModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	notifications := make([]*entities.Notification, len(notificationModels))
	for i, notification := range notificationModels {
		notifications[i] = mapper.FromNotificationModel(notification)
	}

	return notifications, &pagingResponse, nil
}
