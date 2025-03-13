package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UserReportEntity struct {
	ReportID       uuid.UUID
	Report         *ReportEntity
	ReportedUserId uuid.UUID
	ReportedUser   *UserForReport
}

func (ur *UserReportEntity) ValidateUserReport() error {
	return validation.ValidateStruct(ur,
		validation.Field(&ur.ReportID, validation.Required),
		validation.Field(&ur.ReportedUserId, validation.Required),
	)
}

func NewUserReport(
	reportReason string,
	reportType consts.ReportType,
	userId,
	reportedUserId uuid.UUID,
) (*UserReportEntity, error) {
	newUserReport, err := NewReport(reportReason, reportType, userId)
	if err != nil {
		return nil, err
	}

	ure := &UserReportEntity{
		ReportID:       newUserReport.ID,
		Report:         newUserReport,
		ReportedUserId: reportedUserId,
	}

	if err := ure.ValidateUserReport(); err != nil {
		return nil, err
	}

	return ure, nil
}
