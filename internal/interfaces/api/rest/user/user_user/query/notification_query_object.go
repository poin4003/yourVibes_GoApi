package query

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
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

func (req *NotificationQueryObject) ToGetManyNotificationQuery(
	userId uuid.UUID,
) (*query.GetManyNotificationQuery, error) {
	return &query.GetManyNotificationQuery{
		UserId:           userId,
		From:             req.From,
		NotificationType: req.NotificationType,
		CreatedAt:        req.CreatedAt,
		SortBy:           req.SortBy,
		IsDescending:     req.IsDescending,
		Limit:            req.Limit,
		Page:             req.Page,
	}, nil
}
