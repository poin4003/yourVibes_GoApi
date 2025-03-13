package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToUserReportModel(userReport *reportEntity.UserReportEntity) *models.UserReport {
	ur := &models.UserReport{
		Report: ToReportModel(userReport.Report),
	}

	ur.ReportID = userReport.ReportID
	ur.ReportedUserId = userReport.ReportedUserId

	return ur
}

func FromUserReportModel(ur *models.UserReport) *reportEntity.UserReportEntity {
	if ur == nil {
		return nil
	}

	var report = &reportEntity.ReportEntity{
		ID:        ur.ReportID,
		UserId:    ur.Report.UserId,
		AdminId:   ur.Report.AdminId,
		User:      *FromUserModel(ur.Report.User),
		Admin:     FromAdminModel(ur.Report.Admin),
		Reason:    ur.Report.Reason,
		Status:    ur.Report.Status,
		CreatedAt: ur.Report.CreatedAt,
		UpdatedAt: ur.Report.UpdatedAt,
	}

	var userReport = &reportEntity.UserReportEntity{
		Report:       report,
		ReportedUser: FromUserModel(ur.ReportedUser),
	}

	userReport.ReportID = ur.ReportID
	userReport.ReportedUserId = ur.ReportedUserId

	return userReport
}
