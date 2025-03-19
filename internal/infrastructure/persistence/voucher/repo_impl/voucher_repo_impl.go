package repo_impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/voucher/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"gorm.io/gorm"
)

type rVoucher struct {
	db *gorm.DB
}

func NewVoucherRepositoryImplement(db *gorm.DB) *rVoucher {
	return &rVoucher{db: db}
}

func (r *rVoucher) GetVoucherById(
	ctx context.Context,
	voucherId uuid.UUID,
) (*entities.VoucherEntity, error) {
	var voucherModel models.Voucher
	if err := r.db.WithContext(ctx).
		Preload("Admin").
		First(&voucherModel, voucherId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return mapper.FromVoucherModel(&voucherModel), nil
}

func (r *rVoucher) GetVoucherByCode(
	ctx context.Context,
	voucherCode string,
) (*entities.VoucherEntity, error) {
	var voucherModel models.Voucher
	if err := r.db.WithContext(ctx).
		First(&voucherModel).
		Where("code = ?", voucherCode).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return mapper.FromVoucherModel(&voucherModel), nil
}

func (r *rVoucher) CreateVoucher(
	ctx context.Context,
	voucher *entities.VoucherEntity,
) error {
	voucherModel := mapper.ToVoucherModel(voucher)
	// 1. Check exists
	voucherExists, _ := r.checkVoucherExistsByCode(ctx, voucherModel.Code)
	if voucherExists {
		return response.NewCustomError(response.ErrDataHasAlreadyExist, "voucher already exists")
	}

	// 2. Create voucher
	if err := r.db.WithContext(ctx).
		Create(&voucherModel).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rVoucher) RedeemVoucher(
	ctx context.Context,
	voucherCode string,
) (*entities.VoucherEntity, error) {
	voucherModel := &models.Voucher{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. find voucher
		if err := tx.WithContext(ctx).
			First(voucherModel).
			Where("code = ?", voucherCode).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Check max uses of voucher
		if voucherModel.MaxUses <= 0 {
			// Delete voucher if max use = 0
			if err := tx.WithContext(ctx).
				Model(voucherModel).
				Where("id = ?", voucherModel.ID).
				Delete(voucherModel).
				Error; err != nil {
				return response.NewServerFailedError(err.Error())
			}
			return response.NewCustomError(response.ErrVoucherExpired)
		}

		// 3. Update Max uses of voucher
		if err := tx.WithContext(ctx).
			Model(voucherModel).
			Where("id = ?", voucherModel.ID).
			Update("max_uses", gorm.Expr("max_uses - ?", 1)).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromVoucherModel(voucherModel), nil
}

func (r *rVoucher) checkVoucherExistsByCode(
	ctx context.Context,
	code string,
) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Voucher{}).
		Where("code = ?", code).
		Count(&count).
		Error; err != nil {
	}
	return count > 0, nil
}
