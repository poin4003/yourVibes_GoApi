package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"

	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
)

type (
	IAdvertise interface {
		CreateAdvertise(ctx context.Context, command *command.CreateAdvertiseCommand) (result *command.CreateAdvertiseResult, err error)
		GetManyAdvertise(ctx context.Context, query *advertise_query.GetManyAdvertiseQuery) (result *advertise_query.GetManyAdvertiseResults, err error)
		GetAdvertise(ctx context.Context, query *advertise_query.GetOneAdvertiseQuery) (result *advertise_query.GetOneAdvertiseResult, err error)
		GetAdvertiseWithStatistic(ctx context.Context, AdvertiseId uuid.UUID) (result *common.AdvertiseForStatisticResult, err error)
		GetShortAdvertiseByUserId(ctx context.Context, query *advertise_query.GetManyAdvertiseByUserId) (result *advertise_query.GetManyAdvertiseResultsByUserId, err error)
	}
	IBill interface {
		ConfirmPayment(ctx context.Context, command *command.ConfirmPaymentCommand) error
	}
)
