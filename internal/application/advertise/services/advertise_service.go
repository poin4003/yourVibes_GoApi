package services

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
)

type (
	IAdvertise interface {
		CreateAdvertise(ctx context.Context, command *command.CreateAdvertiseCommand) (result *command.CreateAdvertiseResult, err error)
		GetManyAdvertise(ctx context.Context, query *advertise_query.GetManyAdvertiseQuery) (result *advertise_query.GetManyAdvertiseResults, err error)
		GetAdvertise(ctx context.Context, query *advertise_query.GetOneAdvertiseQuery) (result *advertise_query.GetOneAdvertiseResult, err error)
	}
	IBill interface {
		ConfirmPayment(ctx context.Context, command *command.ConfirmPaymentCommand) error
	}
)

var (
	localAdvertise IAdvertise
	localBill      IBill
)

func Advertise() IAdvertise {
	if localAdvertise == nil {
		panic("service_implement localAdvertise not found for interface IAdvertise")
	}
	return localAdvertise
}

func Bill() IBill {
	if localBill == nil {
		panic("service_implement localBill not found for interface IBill")
	}

	return localBill
}

func InitAdvertise(i IAdvertise) {
	localAdvertise = i
}

func InitBill(i IBill) {
	localBill = i
}
