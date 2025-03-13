package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToReportModel(report *entities.ReportEntity) *models.Report {
	r := &models.Report{
		UserId:    report.UserId,
		AdminId:   report.AdminId,
		Reason:    report.Reason,
		Type:      report.Type,
		Status:    report.Status,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
	}
	r.ID = report.ID

	return r
}
