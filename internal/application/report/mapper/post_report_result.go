package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewPostReportResult(
	postReport *reportEntity.PostReportEntity,
) *common.PostReportResult {
	if postReport == nil {
		return nil
	}

	var postReportResult = &common.PostReportResult{
		AdminId:        postReport.Report.AdminId,
		UserId:         postReport.Report.UserId,
		ReportedPostId: postReport.ReportedPostId,
		User:           NewUserResult(&postReport.Report.User),
		ReportedPost:   NewPostResult(postReport.ReportedPost),
		Admin:          NewAdminResult(postReport.Report.Admin),
		Reason:         postReport.Report.Reason,
		Status:         postReport.Report.Status,
		CreatedAt:      postReport.Report.CreatedAt,
		UpdatedAt:      postReport.Report.UpdatedAt,
	}
	postReportResult.ReportId = postReport.ReportID

	return postReportResult
}

func NewPostReportShortVerResult(
	postReport *reportEntity.PostReportEntity,
) *common.PostReportShortVerResult {
	if postReport == nil {
		return nil
	}

	var adminEmail *string
	if postReport.Report.Admin != nil {
		adminEmail = &postReport.Report.Admin.Email
	}

	var postReportResult = &common.PostReportShortVerResult{
		UserId:         postReport.Report.UserId,
		AdminId:        postReport.Report.AdminId,
		UserEmail:      postReport.Report.User.Email,
		ReportedPostId: postReport.ReportedPostId,
		AdminEmail:     adminEmail,
		Reason:         postReport.Report.Reason,
		Status:         postReport.Report.Status,
		CreatedAt:      postReport.Report.CreatedAt,
		UpdatedAt:      postReport.Report.UpdatedAt,
	}
	postReportResult.ReportId = postReport.ReportID

	return postReportResult
}
