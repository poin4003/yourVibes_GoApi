package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToVoucherModel(voucher *entities.VoucherEntity) *models.Voucher {
	v := &models.Voucher{
		ID:          voucher.ID,
		AdminId:     voucher.AdminId,
		Name:        voucher.Name,
		Description: voucher.Description,
		Type:        voucher.Type,
		Value:       voucher.Value,
		Code:        voucher.Code,
		MaxUses:     voucher.MaxUses,
		Status:      voucher.Status,
		CreatedAt:   voucher.CreatedAt,
		UpdatedAt:   voucher.UpdatedAt,
	}

	return v
}

func FromVoucherModel(voucherModel *models.Voucher) *entities.VoucherEntity {
	if voucherModel == nil {
		return nil
	}

	return &entities.VoucherEntity{
		ID:          voucherModel.ID,
		AdminId:     voucherModel.AdminId,
		Admin:       FromAdminModel(voucherModel.Admin),
		Name:        voucherModel.Name,
		Description: voucherModel.Description,
		Type:        voucherModel.Type,
		Value:       voucherModel.Value,
		Code:        voucherModel.Code,
		MaxUses:     voucherModel.MaxUses,
		Status:      voucherModel.Status,
		CreatedAt:   voucherModel.CreatedAt,
		UpdatedAt:   voucherModel.UpdatedAt,
	}
}
