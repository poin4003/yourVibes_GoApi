package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type Setting struct {
	ID        uint
	UserId    uuid.UUID
	Language  consts.Language
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SettingUpdate struct {
	Language  *consts.Language
	Status    *bool
	UpdatedAt *time.Time
}

func (s *Setting) ValidateSetting() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.UserId, validation.Required),
		validation.Field(&s.Language, validation.Required, validation.In(consts.Languages...)),
		validation.Field(&s.Status, validation.Required),
		validation.Field(&s.CreatedAt, validation.Required),
		validation.Field(&s.UpdatedAt, validation.Required, validation.Min(s.CreatedAt)),
	)
}

func (s *SettingUpdate) ValidateSettingUpdate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Language, validation.In(consts.Languages...)),
	)
}

func NewSetting(
	userId uuid.UUID,
	language consts.Language,
) (*Setting, error) {
	setting := &Setting{
		UserId:    userId,
		Language:  language,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := setting.ValidateSetting(); err != nil {
		return nil, err
	}

	return setting, nil
}

func NewSettingUpdate(
	updateData *SettingUpdate,
) (*SettingUpdate, error) {
	settingUpdate := &SettingUpdate{
		Language:  updateData.Language,
		Status:    updateData.Status,
		UpdatedAt: updateData.UpdatedAt,
	}

	if err := settingUpdate.ValidateSettingUpdate(); err != nil {
		return nil, err
	}

	return settingUpdate, nil
}
