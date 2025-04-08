package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/statistic/entities"
)

type (
	IStatisticRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.StatisticEntity, error)
		UpsertLatestStatistic(ctx context.Context, postId uuid.UUID, entity *entities.StatisticEntity) error
	}
)
