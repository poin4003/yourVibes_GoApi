package query

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/common"
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
	Notifications  []*common.NotificationResultForInterface
	PagingResponse *response.PagingResponse
}
