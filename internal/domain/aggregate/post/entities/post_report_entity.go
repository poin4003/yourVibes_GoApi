package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type PostReport struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	AdminId        *uuid.UUID
	User           *UserForReport
	ReportedPost   *PostForReport
	Admin          *Admin
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PostReportUpdate struct {
	AdminId *uuid.UUID
	Status  *bool
}

func (pr *PostReport) ValidatePostReport() error {
	return validation.ValidateStruct(pr,
		validation.Field(&pr.UserId, validation.Required),
		validation.Field(&pr.ReportedPostId, validation.Required),
		validation.Field(&pr.Reason, validation.Required, validation.Length(10, 255)),
	)
}

func NewPostReport(
	userId uuid.UUID,
	reportedPostId uuid.UUID,
	reason string,
) (*PostReport, error) {
	newPostReport := &PostReport{
		UserId:         userId,
		ReportedPostId: reportedPostId,
		AdminId:        &uuid.Nil,
		Reason:         reason,
		Status:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := newPostReport.ValidatePostReport(); err != nil {
		return nil, err
	}

	return newPostReport, nil
}
