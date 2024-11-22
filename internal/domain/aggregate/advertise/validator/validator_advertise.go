package validator

import (
	"fmt"
	advertise_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

type ValidateAdvertise struct {
	advertise_entity.Advertise
	isValidated bool
}

func (va *ValidateAdvertise) Valid() bool {
	return va.isValidated
}

func NewValidateAdvertise(
	advertise *advertise_entity.Advertise,
) (*ValidateAdvertise, error) {
	if advertise == nil {
		return nil, fmt.Errorf("advertise is nil")
	}

	if err := advertise.Validate(); err != nil {
		return nil, err
	}

	return &ValidateAdvertise{
		Advertise:   *advertise,
		isValidated: true,
	}, nil
}
