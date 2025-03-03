package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type GetManyNotificationQuery struct {
	UserId           uuid.UUID
	From             string
	NotificationType string
	CreatedAt        time.Time
	SortBy           string
	IsDescending     bool
	Limit            int
	Page             int
}

type GetManyNotificationQueryResult struct {
	Notifications  []*common.NotificationResult
	PagingResponse *response.PagingResponse
}
