package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"time"
)

type Media struct {
	ID        uint
	PostId    uuid.UUID
	MediaUrl  string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MediaUpdate struct {
	MediaUrl  *string
	Status    *bool
	UpdatedAt *time.Time
}

func (m *Media) ValidateMedia() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.MediaUrl, validation.Required, is.URL),
	)
}

func (m *MediaUpdate) ValidateMediaUpdate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.MediaUrl, is.URL),
	)
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

	if err := media.ValidateMedia(); err != nil {
		return nil, err
	}

	return media, nil
}
