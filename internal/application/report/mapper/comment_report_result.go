package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewCommentReportResult(
	commentReport *reportEntity.CommentReportEntity,
) *common.CommentReportResult {
	if commentReport == nil {
		return nil
	}

	var commentReportResult = &common.CommentReportResult{
		AdminId:           commentReport.Report.AdminId,
		UserId:            commentReport.Report.UserId,
		ReportedCommentId: commentReport.ReportedCommentId,
		User:              NewUserResult(&commentReport.Report.User),
		ReportedComment:   NewCommentResult(commentReport.ReportedComment),
		Post:              NewPostResult(commentReport.Post),
		Admin:             NewAdminResult(commentReport.Report.Admin),
		Reason:            commentReport.Report.Reason,
		Status:            commentReport.Report.Status,
		CreatedAt:         commentReport.Report.CreatedAt,
		UpdatedAt:         commentReport.Report.UpdatedAt,
	}
	commentReportResult.ReportId = commentReport.ReportID

	return commentReportResult
}

func NewCommentReportShortVerResult(
	commentReport *reportEntity.CommentReportEntity,
) *common.CommentReportShortVerResult {
	if commentReport == nil {
		return nil
	}

	var adminEmail *string
	if commentReport.Report.Admin != nil {
		adminEmail = &commentReport.Report.Admin.Email
	}

	var commentReportResult = &common.CommentReportShortVerResult{
		UserId:            commentReport.Report.UserId,
		AdminId:           commentReport.Report.AdminId,
		Reason:            commentReport.Report.Reason,
		UserEmail:         commentReport.Report.User.Email,
		ReportedCommentId: commentReport.ReportedCommentId,
		AdminEmail:        adminEmail,
		Status:            commentReport.Report.Status,
		CreatedAt:         commentReport.Report.CreatedAt,
		UpdatedAt:         commentReport.Report.UpdatedAt,
	}
	commentReportResult.ReportId = commentReport.ReportID

	return commentReportResult
}
