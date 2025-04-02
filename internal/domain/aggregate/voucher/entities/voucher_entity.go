package entities

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/random"
)

type VoucherEntity struct {
	ID          uuid.UUID
	AdminId     *uuid.UUID
	Admin       *Admin
	Name        string
	Description string
	Type        consts.VoucherType
	Value       int
	Code        string
	MaxUses     int
	Status      consts.VoucherStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ve *VoucherEntity) ValidateVoucherEntity() error {
	return validation.ValidateStruct(ve,
		validation.Field(&ve.Name, validation.Required, validation.RuneLength(2, 50)),
		validation.Field(&ve.Description, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&ve.Type, validation.Required, validation.In(consts.VoucherTypes...)),
		validation.Field(&ve.Value, validation.Required, validation.By(func(value interface{}) error {
			v := value.(int)
			if ve.Type == consts.PERCENTAGE {
				if v < 1 || v > 100 {
					return fmt.Errorf("value must be between 1 and 100 if type is percentage")
				}
			} else {
				if v < 0 {
					return fmt.Errorf("value must be greater or equal to 0")
				}
			}
			return nil
		})),
		validation.Field(&ve.Code, validation.Required, validation.RuneLength(1, 30)),
		validation.Field(&ve.MaxUses, validation.Required, validation.Min(0)),
	)
}

func NewVoucherBySystem(
	name, description string,
	maxUses, value int,
	voucherType consts.VoucherType,
) (*VoucherEntity, error) {
	voucher := &VoucherEntity{
		ID:          uuid.New(),
		AdminId:     nil,
		Name:        name,
		Description: description,
		Type:        voucherType,
		Value:       value,
		Code:        random.GenerateVoucherCode(name),
		MaxUses:     maxUses,
		Status:      consts.VOUCHER_ACTIVE,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := voucher.ValidateVoucherEntity(); err != nil {
		return nil, err
	}

	return voucher, nil
}

func NewVoucherByAdmin(
	name, description, code string,
	maxUses, value int,
	voucherType consts.VoucherType,
	adminId uuid.UUID,
) (*VoucherEntity, error) {
	voucher := &VoucherEntity{
		ID:          uuid.New(),
		AdminId:     &adminId,
		Name:        name,
		Description: description,
		Type:        voucherType,
		Value:       value,
		Code:        code,
		MaxUses:     maxUses,
		Status:      consts.VOUCHER_ACTIVE,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := voucher.ValidateVoucherEntity(); err != nil {
		return nil, err
	}

	return voucher, nil
}
