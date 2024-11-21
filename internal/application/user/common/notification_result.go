package common

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type NotificationResult struct {
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
