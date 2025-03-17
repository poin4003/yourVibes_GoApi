package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type NotificationResult struct {
	From             string                  `json:"from"`
	FromUrl          string                  `json:"from_url"`
	UserID           string                  `json:"user_id"`
	NotificationType consts.NotificationType `json:"notification_type"`
	ContentID        string                  `json:"content_id"`
	Content          string                  `json:"content"`
	Status           bool                    `json:"status"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

type NotificationResultForInterface struct {
	ID               uint
	From             string
	FromUrl          string
	UserId           uuid.UUID
	User             *UserShortVerResult
	NotificationType consts.NotificationType
	ContentId        string
	Content          string
	Status           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
