package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
)

type (
	IVoucherRepository interface {
		GetVoucherById(ctx context.Context, voucherId uuid.UUID) (*entities.VoucherEntity, error)
		GetVoucherByCode(ctx context.Context, voucherCode string) (*entities.VoucherEntity, error)
		CreateVoucher(ctx context.Context, voucher *entities.VoucherEntity) error
		RedeemVoucher(ctx context.Context, voucherCode string) (*entities.VoucherEntity, error)
	}
)

var (
	localVoucher IVoucherRepository
)

func Voucher() IVoucherRepository {
	if localVoucher == nil {
		panic("repository implementation of Voucher not found")
	}

	return localVoucher
}

func InitVoucherRepository(i IVoucherRepository) {
	localVoucher = i
}
