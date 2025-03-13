package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToCommentReportModel(commentReport *reportEntity.CommentReportEntity) *models.CommentReport {
	cr := &models.CommentReport{
		Report: ToReportModel(commentReport.Report),
	}

	cr.ReportID = commentReport.ReportID
	cr.ReportedCommentId = commentReport.ReportedCommentId

	return cr
}

func FromCommentReportModel(cr *models.CommentReport) *reportEntity.CommentReportEntity {
	if cr == nil {
		return nil
	}

	var report = &reportEntity.ReportEntity{
		ID:        cr.ReportID,
		UserId:    cr.Report.UserId,
		AdminId:   cr.Report.AdminId,
		User:      *FromUserModel(cr.Report.User),
		Admin:     FromAdminModel(cr.Report.Admin),
		Reason:    cr.Report.Reason,
		Status:    cr.Report.Status,
		CreatedAt: cr.Report.CreatedAt,
		UpdatedAt: cr.Report.UpdatedAt,
	}

	var commentReport = &reportEntity.CommentReportEntity{
		Report:          report,
		ReportedComment: FromCommentModel(cr.ReportedComment),
		Post:            nil,
	}

	if cr.ReportedComment != nil && cr.ReportedComment.Post != nil {
		commentReport.Post = FromPostModel(cr.ReportedComment.Post)
	}

	commentReport.ReportID = cr.ReportID
	commentReport.ReportedCommentId = cr.ReportedCommentId

	return commentReport
}
