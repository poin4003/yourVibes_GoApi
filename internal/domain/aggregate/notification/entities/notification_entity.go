package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type Notification struct {
	ID               uint                    `validated:"omitempty"`
	From             string                  `validated:"required,min=2"`
	FromUrl          string                  `validated:"required,url"`
	UserId           uuid.UUID               `validated:"required,uuid4"`
	User             *User                   `validated:"required"`
	NotificationType consts.NotificationType `validated:"required,notification_type"`
	ContentId        string                  `validated:"required,min=2"`
	Content          string                  `validated:"required,min=2"`
	Status           bool                    `validated:"required"`
	CreatedAt        time.Time               `validated:"required"`
	UpdatedAt        time.Time               `validated:"required,gtefield=CreatedAt"`
}

type NotificationUpdate struct {
	From             *string                  `validated:"omitempty,min=2"`
	FromUrl          *string                  `validated:"omitempty,url"`
	NotificationType *consts.NotificationType `validated:"omitempty,notification_type"`
	ContentId        *string                  `validated:"omitempty,min=2"`
	Content          *string                  `validated:"omitempty,min=2"`
	Status           *bool                    `validated:"omitempty"`
	UpdatedAt        *time.Time               `validated:"omitempty,gtefield=CreatedAt"`
}

func (n *Notification) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("notification_type", func(fl validator.FieldLevel) bool {
		notificationType := consts.NotificationType(fl.Field().String())
		switch notificationType {
		case consts.NEW_POST,
			consts.NEW_COMMENT,
			consts.LIKE_POST,
			consts.LIKE_COMMENT,
			consts.NEW_SHARE,
			consts.FRIEND_REQUEST,
			consts.ACCEPT_FRIEND_REQUEST:
			return true
		default:
			return false
		}
	})
	return validate.Struct(n)
}

func (n *NotificationUpdate) ValidateNotificationUpdate() error {
	validate := validator.New()
	validate.RegisterValidation("notification_type", func(fl validator.FieldLevel) bool {
		notificationType := consts.NotificationType(fl.Field().String())
		switch notificationType {
		case consts.NEW_POST,
			consts.NEW_COMMENT,
			consts.LIKE_POST,
			consts.LIKE_COMMENT,
			consts.NEW_SHARE,
			consts.FRIEND_REQUEST,
			consts.ACCEPT_FRIEND_REQUEST:
			return true
		default:
			return false
		}
	})
	return validate.Struct(n)
}

func NewNotification(
	from string,
	fromUrl string,
	userId uuid.UUID,
	notificationType consts.NotificationType,
	contentId string,
	content string,
) (*Notification, error) {
	notification := &Notification{
		From:             from,
		FromUrl:          fromUrl,
		UserId:           userId,
		NotificationType: notificationType,
		ContentId:        contentId,
		Content:          content,
		Status:           true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := notification.Validate(); err != nil {
		return nil, err
	}

	return notification, nil
}

func NewNotificationUpdate(
	updateData *NotificationUpdate,
) (*NotificationUpdate, error) {
	notificationUpdate := &NotificationUpdate{
		From:             updateData.From,
		FromUrl:          updateData.FromUrl,
		NotificationType: updateData.NotificationType,
		ContentId:        updateData.ContentId,
		Content:          updateData.Content,
		Status:           updateData.Status,
		UpdatedAt:        updateData.UpdatedAt,
	}

	if err := notificationUpdate.ValidateNotificationUpdate(); err != nil {
		return nil, err
	}

	return notificationUpdate, nil
}
