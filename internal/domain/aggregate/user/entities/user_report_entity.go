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
	User           *UserForReport
	ReportedUser   *UserForReport
	Admin          *Admin
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserReportUpdate struct {
	AdminId *uuid.UUID
	Status  *bool
}

func (ur *UserReport) ValidateUserReport() error {
	return validation.ValidateStruct(ur,
		validation.Field(&ur.UserId, validation.Required),
		validation.Field(&ur.ReportedUserId, validation.Required),
		validation.Field(&ur.Reason, validation.Required, validation.Length(10, 255)),
	)
}

func NewUserReport(
	userId uuid.UUID,
	reportedUserId uuid.UUID,
	reason string,
) (*UserReport, error) {
	newUserReport := &UserReport{
		UserId:         userId,
		ReportedUserId: reportedUserId,
		AdminId:        uuid.Nil,
		Reason:         reason,
		Status:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := newUserReport.ValidateUserReport(); err != nil {
		return nil, err
	}

	return newUserReport, nil
}
