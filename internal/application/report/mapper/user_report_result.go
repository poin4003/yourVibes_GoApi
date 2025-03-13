package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewUserReportResult(
	userReport *reportEntity.UserReportEntity,
) *common.UserReportResult {
	if userReport == nil {
		return nil
	}

	var userReportResult = &common.UserReportResult{
		AdminId:        userReport.Report.AdminId,
		UserId:         userReport.Report.UserId,
		ReportedUserId: userReport.ReportedUserId,
		User:           NewUserResult(&userReport.Report.User),
		ReportedUser:   NewUserResult(userReport.ReportedUser),
		Admin:          NewAdminResult(userReport.Report.Admin),
		Reason:         userReport.Report.Reason,
		Status:         userReport.Report.Status,
		CreatedAt:      userReport.Report.CreatedAt,
		UpdatedAt:      userReport.Report.UpdatedAt,
	}
	userReportResult.ReportId = userReport.ReportID

	return userReportResult
}

func NewUserReportShortVerResult(
	userReport *reportEntity.UserReportEntity,
) *common.UserReportShortVerResult {
	if userReport == nil {
		return nil
	}

	var adminEmail *string
	if userReport.Report.Admin != nil {
		adminEmail = &userReport.Report.Admin.Email
	}

	var userReportResult = &common.UserReportShortVerResult{
		UserId:            userReport.Report.UserId,
		AdminId:           userReport.Report.AdminId,
		UserEmail:         userReport.Report.User.Email,
		ReportedUserId:    userReport.ReportedUserId,
		ReportedUserEmail: userReport.ReportedUser.Email,
		AdminEmail:        adminEmail,
		Reason:            userReport.Report.Reason,
		Status:            userReport.Report.Status,
		CreatedAt:         userReport.Report.CreatedAt,
		UpdatedAt:         userReport.Report.UpdatedAt,
	}
	userReportResult.ReportId = userReport.ReportID

	return userReportResult
}
