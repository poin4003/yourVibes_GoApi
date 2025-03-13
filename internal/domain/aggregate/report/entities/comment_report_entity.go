package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CommentReportEntity struct {
	ReportID          uuid.UUID
	Report            *ReportEntity
	ReportedCommentId uuid.UUID
	ReportedComment   *CommentForReport
	Post              *PostForReport
}

func (cr *CommentReportEntity) ValidateCommentReport() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.ReportID, validation.Required),
		validation.Field(&cr.ReportedCommentId, validation.Required),
	)
}

func NewCommentReport(
	reportReason string,
	reportType consts.ReportType,
	userId,
	reportedCommentId uuid.UUID,
) (*CommentReportEntity, error) {
	newReport, err := NewReport(reportReason, reportType, userId)
	if err != nil {
		return nil, err
	}

	cre := &CommentReportEntity{
		ReportID:          newReport.ID,
		Report:            newReport,
		ReportedCommentId: reportedCommentId,
	}

	if err := cre.ValidateCommentReport(); err != nil {
		return nil, err
	}

	return cre, nil
}
