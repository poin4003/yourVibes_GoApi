package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Media struct {
	ID        uint      `validate:"omitempty,uuid4"`
	PostId    uuid.UUID `validate:"required,uuid4"`
	MediaUrl  string    `validate:"required,url"`
	Status    bool      `validate:"omitempty"`
	CreatedAt time.Time `validate:"omitempty"`
	UpdatedAt time.Time `validate:"omitempty,gtefield=CreatedAt"`
}

type MediaUpdate struct {
	MediaUrl  *string    `validate:"omitempty,url"`
	Status    *bool      `validate:"omitempty,required"`
	UpdatedAt *time.Time `validate:"omitempty,gtefield=CreatedAt"`
}

func (m *Media) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *MediaUpdate) ValidateMediaUpdate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func NewMedia(
	PostId uuid.UUID,
	MediaUrl string,
) (*Media, error) {
	media := &Media{
		PostId:    PostId,
		MediaUrl:  MediaUrl,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := media.Validate(); err != nil {
		return nil, err
	}

	return media, nil
}
