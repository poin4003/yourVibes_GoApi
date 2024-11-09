package common

import (
	"github.com/google/uuid"
	"time"
)

type NotificationResult struct {
	ID               uint               `json:"id"`
	From             string             `json:"from"`
	FromUrl          string             `json:"from_url"`
	UserId           uuid.UUID          `json:"user_id"`
	User             UserShortVerResult `json:"user"`
	NotificationType string             `json:"notification_type"`
	ContentId        string             `json:"content_id"`
	Content          string             `json:"content"`
	Status           bool               `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
