package services

import "context"

type (
	IAdvertise interface {
		CreateAdvertise(ctx context.Context)
		GetAdvertise(ctx context.Context)
	}
	IBill interface {
		ConfirmPayment(ctx context.Context)
	}
)
