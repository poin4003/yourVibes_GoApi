package repo_impl

import (
	"context"
	"errors"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/notification/mapper"
	"gorm.io/gorm"
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		First(&notificationModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
	for _, notification := range notificationEntities {
		notificationModels = append(notificationModels, mapper.ToNotificationModel(notification))
	}

	err := r.db.WithContext(ctx).Create(&notificationModels).Error
	if err != nil {
		return nil, err
	}

	var notificationEntityList []*entities.Notification
	for _, notificationEntity := range notificationModels {
		notificationEntityList = append(notificationEntityList, mapper.FromNotificationModel(notificationEntity))
	}

	return notificationEntityList, nil
}

func (r *rNotification) UpdateOne(
	ctx context.Context,
	notificationId uint,
	updateData *entities.NotificationUpdate,
) (*entities.Notification, error) {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no field to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", notificationId).
		Updates(&updates).
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
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
		sortColumn := ""
		switch query.SortBy {
		case "id":
			sortColumn = "id"
		case "from":
			sortColumn = "from"
		case "notification_type":
			sortColumn = "notification_type"
		case "created_at":
			sortColumn = "created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
		}
	}

	err := db.Where("user_id = ?", query.UserId).Count(&total).Error
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Find(&notificationModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var notificationEntities []*entities.Notification
	for _, notification := range notificationModels {
		notificationEntities = append(notificationEntities, mapper.FromNotificationModel(notification))
	}

	return notificationEntities, &pagingResponse, nil
}
