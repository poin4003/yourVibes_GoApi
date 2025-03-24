package implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/command"
	statisticEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/statistic/entities"
	repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type sStatistic struct {
	statisticRepo repo.IStatisticRepository
}

func NewStatisticImplement(
	statisticRepo repo.IStatisticRepository,
) *sStatistic {
	return &sStatistic{
		statisticRepo: statisticRepo,
	}
}

func (s *sStatistic) UpsertStatistic(
	ctx context.Context,
	postId uuid.UUID,
	command *command.UpsertStatisticCommand,
) error {
	statistic, err := statisticEntity.NewStatisticEntity(
		postId,
		command.Reach,
		command.Clicks,
		command.Impression,
	)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if err = s.statisticRepo.UpsertLatestStatistic(ctx, postId, statistic); err != nil {
		return err
	}
	return nil
}
