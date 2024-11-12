package entities

import (
	"github.com/google/uuid"
	"time"
	"github.com/go-playground/validator/v10"
)

type Media struct {
	ID        uint           `validate:"required,uuid4"`
	PostId    uuid.UUID      `validate:"required,uuid4"`
	MediaUrl  string         `validate:"required,url"`
	Status    bool           `validate:"required"`
	CreatedAt time.Time      `validate:"required"`
	UpdatedAt time.Time      `validate:"required,gtefield=CreatedAt"`
}

func (m *Media) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}