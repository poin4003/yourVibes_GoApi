package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type CommentReport struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	AdminId           *uuid.UUID
	User              *UserForReport
	ReportedComment   *CommentForReport
	Post              *PostForReport
	Admin             *Admin
	Reason            string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CommentReportUpdate struct {
	AdminId *uuid.UUID
	Status  *bool
}

func (cr *CommentReport) ValidateCommentReport() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.UserId, validation.Required),
		validation.Field(&cr.ReportedCommentId, validation.Required),
		validation.Field(&cr.Reason, validation.Required, validation.Length(2, 255)),
	)
}

func NewCommentReport(
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
	reason string,
) (*CommentReport, error) {
	newCommentReport := &CommentReport{
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
		AdminId:           &uuid.Nil,
		Reason:            reason,
		Status:            false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := newCommentReport.ValidateCommentReport(); err != nil {
		return nil, err
	}

	return newCommentReport, nil
}
