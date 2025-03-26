package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/statistic/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToStatisticModel(statistic *entities.StatisticEntity) *models.Statistics {
	s := &models.Statistics{
		PostId:     statistic.PostId,
		Reach:      statistic.Reach,
		Clicks:     statistic.Clicks,
		Impression: statistic.Impression,
		Status:     statistic.Status,
		CreatedAt:  statistic.CreatedAt,
		UpdatedAt:  statistic.UpdatedAt,
	}
	s.ID = statistic.ID

	return s
}

func FromStatisticModel(s *models.Statistics) *entities.StatisticEntity {
	if s == nil {
		return nil
	}

	statistic := &entities.StatisticEntity{
		PostId:     s.PostId,
		Reach:      s.Reach,
		Clicks:     s.Clicks,
		Impression: s.Impression,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
	statistic.ID = s.ID

	return statistic
}
