package entities

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type Notification struct {
	ID               uint
	From             string
	FromUrl          string
	UserId           uuid.UUID
	User             *User
	NotificationType consts.NotificationType
	ContentId        string
	Content          string
	Status           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type NotificationUpdate struct {
	From             *string
	FromUrl          *string
	NotificationType *consts.NotificationType
	ContentId        *string
	Content          *string
	Status           *bool
	UpdatedAt        *time.Time
}

func (n *Notification) ValidateNotification() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.From, validation.Required, validation.Length(2, 255)),
		validation.Field(&n.FromUrl, validation.Required, is.URL),
		validation.Field(&n.UserId, validation.Required),
		validation.Field(&n.User, validation.Required),
		validation.Field(&n.NotificationType, validation.Required, validation.By(validateNotificationType)),
		validation.Field(&n.ContentId, validation.Required, validation.Length(2, 255)),
		validation.Field(&n.Content, validation.Required, validation.Length(2, 255)),
		validation.Field(&n.Status, validation.Required),
		validation.Field(&n.CreatedAt, validation.Required),
		validation.Field(&n.UpdatedAt, validation.Required, validation.Min(n.CreatedAt)),
	)
}

func (n *NotificationUpdate) ValidateNotificationUpdate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.From, validation.Length(2, 0)),
		validation.Field(&n.FromUrl, is.URL),
		validation.Field(&n.NotificationType, validation.By(validateNotificationType)),
		validation.Field(&n.ContentId, validation.Length(2, 0)),
		validation.Field(&n.Content, validation.Length(2, 0)),
	)
}

func validateNotificationType(value interface{}) error {
	notificationType, ok := value.(consts.NotificationType)
	if !ok {
		return fmt.Errorf("invalid notification type")
	}

	switch notificationType {
	case consts.NEW_POST, consts.NEW_COMMENT, consts.LIKE_POST, consts.LIKE_COMMENT,
		consts.NEW_SHARE, consts.FRIEND_REQUEST, consts.ACCEPT_FRIEND_REQUEST:
		return nil
	default:
		return fmt.Errorf("invalid notification type: %v", notificationType)
	}
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
	if err := notification.ValidateNotification(); err != nil {
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
