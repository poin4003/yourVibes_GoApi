package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
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

func (r *rNotification) CreateNotification(
	ctx context.Context,
	notification *model.Notification,
) (*model.Notification, error) {
	res := r.db.WithContext(ctx).Preload("User").Create(notification)

	if res.Error != nil {
		return nil, res.Error
	}

	return notification, nil
}

func (r *rNotification) CreateManyNotification(
	ctx context.Context,
	notifications []*model.Notification,
) ([]*model.Notification, error) {
	err := r.db.WithContext(ctx).Create(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *rNotification) UpdateOneNotification(
	ctx context.Context,
	notificationId uint,
	updateData map[string]interface{},
) (*model.Notification, error) {
	var notification model.Notification

	if err := r.db.WithContext(ctx).First(&notification, notificationId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&notification).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *rNotification) UpdateManyNotification(
	ctx context.Context,
	condition map[string]interface{},
	updateData map[string]interface{},
) error {
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where(condition).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *rNotification) DeleteNotification(
	ctx context.Context,
	notificationId uint,
) (*model.Notification, error) {
	notification := &model.Notification{}
	res := r.db.WithContext(ctx).First(notification, notificationId)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Delete(notification)
	if res.Error != nil {
		return nil, res.Error
	}

	return notification, nil
}

func (r *rNotification) GetOneNotification(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Notification, error) {
	notification := &model.Notification{}

	if res := r.db.WithContext(ctx).Model(notification).Where(query, args...).Preload("User").First(&notification); res.Error != nil {
		return nil, res.Error
	}

	return notification, nil
}

func (r *rNotification) GetManyNotification(
	ctx context.Context,
	userId uuid.UUID,
	query *query_object.NotificationQueryObject,
) ([]*model.Notification, *response.PagingResponse, error) {
	var notifications []*model.Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Notification{})

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
		Where("user_id=?", userId).
		Preload("User").
		Find(&notifications).Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return notifications, &pagingResponse, nil
}
