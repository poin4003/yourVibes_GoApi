package validator

import (
	"fmt"
	admin_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
)

type ValidateAdmin struct {
	admin_entity.Admin
	isValidated bool
}

func (vad *ValidateAdmin) Valid() bool {
	return vad.isValidated
}

func NewValidateAdmin(admin *admin_entity.Admin) (*ValidateAdmin, error) {
	if admin == nil {
		return nil, fmt.Errorf("NewValidatedAdmin: admin is nil")
	}

	if err := admin.ValidateAdmin(); err != nil {
		return nil, err
	}

	return &ValidateAdmin{*admin, false}, nil
}
