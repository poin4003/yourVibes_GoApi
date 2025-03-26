package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"time"
)

type StatisticDto struct {
	PostId          uuid.UUID `json:"post_id"`
	Reach           int       `json:"reach"`
	Clicks          int       `json:"clicks"`
	Impression      int       `json:"impression"`
	AggregationDate time.Time `json:"aggregation_date"`
}

func ToStatisticDto(statisticResult *common.StatisticResult) *StatisticDto {
	if statisticResult == nil {
		return nil
	}

	return &StatisticDto{
		PostId:          statisticResult.PostId,
		Reach:           statisticResult.Reach,
		Clicks:          statisticResult.Clicks,
		Impression:      statisticResult.Impression,
		AggregationDate: statisticResult.AggregationDate,
	}
}
