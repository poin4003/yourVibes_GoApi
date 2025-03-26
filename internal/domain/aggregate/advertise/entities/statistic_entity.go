package entities

import (
	"github.com/google/uuid"
	"time"
)

type StatisticEntity struct {
	PostId          uuid.UUID
	Reach           int
	Clicks          int
	Impression      int
	AggregationDate time.Time
}
