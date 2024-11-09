package validator

import user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"

type ValidatedSetting struct {
	user_entity.Setting
	isValidated bool
}

func (vs *ValidatedSetting) IsValid() bool {
	return vs.isValidated
}

func NewValidatedSetting(setting *user_entity.Setting) (*ValidatedSetting, error) {
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	return &ValidatedSetting{
		Setting:     *setting,
		isValidated: true,
	}, nil
}
