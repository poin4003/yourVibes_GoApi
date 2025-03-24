package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	statisticEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

func NewStatisticResult(
	statistic *statisticEntity.StatisticEntity,
) *common.StatisticResult {
	if statistic == nil {
		return nil
	}

	return &common.StatisticResult{
		PostId:          statistic.PostId,
		Reach:           statistic.Reach,
		Clicks:          statistic.Clicks,
		Impression:      statistic.Impression,
		AggregationDate: statistic.AggregationDate,
	}
}
