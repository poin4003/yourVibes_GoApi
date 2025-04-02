package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ReportEntity struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	AdminId   *uuid.UUID
	User      UserForReport
	Admin     *Admin
	Reason    string
	Type      consts.ReportType
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (re *ReportEntity) ValidateReportEntity() error {
	return validation.ValidateStruct(re,
		validation.Field(&re.UserId, validation.Required),
		validation.Field(&re.Reason, validation.Required, validation.RuneLength(2, 255)),
		validation.Field(&re.Type, validation.In(consts.ReportTypes...)),
	)
}

func NewReport(
	reportReason string,
	reportType consts.ReportType,
	userId uuid.UUID,
) (*ReportEntity, error) {
	r := &ReportEntity{
		ID:        uuid.New(),
		UserId:    userId,
		AdminId:   nil,
		Reason:    reportReason,
		Type:      reportType,
		Status:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.ValidateReportEntity(); err != nil {
		return nil, err
	}

	return r, nil
}
