package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/query"
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

func ValidateNotificationQueryObject(input interface{}) error {
	query, ok := input.(*NotificationQueryObject)
	if !ok {
		return fmt.Errorf("validateNotificationQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
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
