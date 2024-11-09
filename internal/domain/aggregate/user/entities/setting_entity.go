package entities

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type Setting struct {
	ID        uint            `validate:"required"`
	UserId    uuid.UUID       `validate:"required,uuid4"`
	Language  consts.Language `validate:"required,oneof=vi en"`
	Status    bool            `validate:"required"`
	CreatedAt time.Time       `validate:"required"`
	UpdatedAt time.Time       `validate:"required,gtefield=CreatedAt"`
}

type SettingUpdate struct {
	Language  *consts.Language `validate:"omitempty,oneof=vi en"`
	Status    *bool            `validate:"omitempty"`
	UpdatedAt *time.Time       `validate:"omitempty,gtefield=CreatedAt"`
}

func (s *Setting) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *SettingUpdate) ValidateSettingUpdate() error {
	validate := validator.New()
	return validate.Struct(s)
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
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	return setting, nil
}

func (s *SettingUpdate) setUpdatedAt() {
	now := time.Now()
	s.UpdatedAt = &now
}

func (s *SettingUpdate) UpdateLanguage(language *consts.Language) error {
	if *language != consts.VI && *language != consts.EN {
		return errors.New("invalid language")
	}
	s.Language = language
	s.setUpdatedAt()
	return s.ValidateSettingUpdate()
}

func (s *SettingUpdate) Activate() error {
	*s.Status = true
	s.setUpdatedAt()
	return s.ValidateSettingUpdate()
}

func (s *SettingUpdate) Deactivate() error {
	*s.Status = false
	s.setUpdatedAt()
	return s.ValidateSettingUpdate()
}
