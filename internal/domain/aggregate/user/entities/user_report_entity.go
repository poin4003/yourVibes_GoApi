package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type UserReport struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
	AdminId        uuid.UUID
	User           User
	ReportedUser   User
	Admin          Admin
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (ur *UserReport) ValidateUserReport() error {
	return validation.ValidateStruct(ur,
		validation.Field(&ur.UserId, validation.Required),
		validation.Field(&ur.ReportedUserId, validation.Required),
		validation.Field(&ur.AdminId, validation.Required),
	)
}
