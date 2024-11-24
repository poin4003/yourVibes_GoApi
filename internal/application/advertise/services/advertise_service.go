package services

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
)

type (
	IAdvertise interface {
		CreateAdvertise(ctx context.Context, command *command.CreateAdvertiseCommand) (result *command.CreateAdvertiseResult, err error)
		GetAdvertise(ctx context.Context, query *query.GetOneAdvertiseQuery) (result *query.GetOneAdvertiseResult, err error)
	}
	IBill interface {
		ConfirmPayment(ctx context.Context, command *command.ConfirmPaymentCommand) (result *command.ConfirmPaymentResult, err error)
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
