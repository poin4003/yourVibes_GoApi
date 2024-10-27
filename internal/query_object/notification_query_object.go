package query_object

import (
	"time"
)

type NotificationQueryObject struct {
	From             string    `json:"from,omitempty"`
	NotificationType string    `json:"notification_type,omitempty"`
	CreatedAt        time.Time `form:"created_at,omitempty"`
	SortBy           string    `form:"sort_by,omitempty"`
	IsDescending     bool      `form:"isDescending,omitempty"`
	Limit            int       `form:"limit,omitempty"`
	Page             int       `form:"page,omitempty"`
}
