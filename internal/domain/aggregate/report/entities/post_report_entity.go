package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostReportEntity struct {
	ReportID       uuid.UUID
	Report         *ReportEntity
	ReportedPostId uuid.UUID
	ReportedPost   *PostForReport
}

func (pr *PostReportEntity) ValidatePostReport() error {
	return validation.ValidateStruct(pr,
		validation.Field(&pr.ReportID, validation.Required),
		validation.Field(&pr.ReportedPostId, validation.Required),
	)
}

func NewPostReport(
	reportReason string,
	reportType consts.ReportType,
	userId,
	reportedPostId uuid.UUID,
) (*PostReportEntity, error) {
	newReport, err := NewReport(reportReason, reportType, userId)
	if err != nil {
		return nil, err
	}

	pre := &PostReportEntity{
		ReportID:       newReport.ID,
		Report:         newReport,
		ReportedPostId: reportedPostId,
	}

	if err := pre.ValidatePostReport(); err != nil {
		return nil, err
	}

	return pre, nil
}
