package command

import (
	"github.com/google/uuid"
	"time"
)

type EventCommand struct {
	PostId    uuid.UUID `json:"post_id"`
	EventType string    `json:"event_type"`
	Count     int       `json:"count"`
	Timestamp time.Time `json:"timestamp"`
}
