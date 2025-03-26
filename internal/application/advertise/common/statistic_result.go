package common

import (
	"github.com/google/uuid"
	"time"
)

type StatisticResult struct {
	PostId          uuid.UUID
	Reach           int
	Clicks          int
	Impression      int
	AggregationDate time.Time
}
