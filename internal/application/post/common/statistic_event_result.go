package common

import (
	"github.com/google/uuid"
	"time"
)

type StatisticEventResult struct {
	PostId    uuid.UUID `json:"post_id"`
	EventType string    `json:"event_type"`
	Count     int       `json:"count"`
	Timestamp time.Time `json:"timestamp"`
}
