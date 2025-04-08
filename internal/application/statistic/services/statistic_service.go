package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/command"
)

type (
	IStatisticMQ interface {
		UpsertStatistic(ctx context.Context, postId uuid.UUID, command *command.UpsertStatisticCommand) error
	}
)
