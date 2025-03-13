package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToPostReportModel(postReport *reportEntity.PostReportEntity) *models.PostReport {
	pr := &models.PostReport{
		Report: ToReportModel(postReport.Report),
	}

	pr.ReportID = postReport.ReportID
	pr.ReportedPostId = postReport.ReportedPostId

	return pr
}

func FromPostReportModel(pr *models.PostReport) *reportEntity.PostReportEntity {
	if pr == nil {
		return nil
	}

	var report = &reportEntity.ReportEntity{
		ID:        pr.ReportID,
		UserId:    pr.Report.UserId,
		AdminId:   pr.Report.AdminId,
		User:      *FromUserModel(pr.Report.User),
		Admin:     FromAdminModel(pr.Report.Admin),
		Reason:    pr.Report.Reason,
		Status:    pr.Report.Status,
		CreatedAt: pr.Report.CreatedAt,
		UpdatedAt: pr.Report.UpdatedAt,
	}

	var postReport = &reportEntity.PostReportEntity{
		Report:       report,
		ReportedPost: FromPostModel(pr.ReportedPost),
	}

	postReport.ReportID = pr.ReportID
	postReport.ReportedPostId = pr.ReportedPostId

	return postReport
}
