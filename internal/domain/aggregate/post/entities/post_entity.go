package entities

import (
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type Post struct {
	ID              uuid.UUID           `validate:"omitempty,uuid4"`
	UserId          uuid.UUID           `validate:"omitempty,uuid4"`
	User            *User                `validate:"required"`
	ParentId        *uuid.UUID          `validate:"omitempty,uuid4"`
	ParentPost      *Post               `validate:"omitempty"`
	Content         string              `validate:"omitempty"`
	LikeCount       int                 `validate:"required"`
	CommentCount    int                 `validate:"required"`
	Privacy         consts.PrivacyLevel `validate:"required,oneof=public private friend_only"`
	Location        string              `validate:"omitempty"`
	IsAdvertisement bool                `validate:"required"`
	Status          bool                `validate:"required"`
	CreatedAt       time.Time           `validate:"required"`
	UpdatedAt       time.Time           `validate:"required,gtefield=CreatedAt"`
	Media           []*Media             `validate:"omitempty"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}