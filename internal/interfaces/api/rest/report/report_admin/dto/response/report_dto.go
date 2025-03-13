package response

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

func ToReportResponse(reportResult *query.ReportQueryResult) interface{} {
	switch reportResult.Type {
	case consts.USER_REPORT:
		return ToUserReportDto(reportResult.UserReport)
	case consts.POST_REPORT:
		return ToPostReportDto(reportResult.PostReport)
	case consts.COMMENT_REPORT:
		return ToCommentReportDto(reportResult.CommentReport)
	default:
		return nil
	}
}

func ToReportShortVerResponse(reportResult *query.ReportQueryListResult) interface{} {
	switch reportResult.Type {
	case consts.USER_REPORT:
		var userReportDtos []*UserReportShortVerDto
		for _, userReportResult := range reportResult.UserReports {
			userReportDtos = append(userReportDtos, ToUserReportShortVerDto(userReportResult))
		}
		return userReportDtos
	case consts.POST_REPORT:
		var postReportDtos []*PostReportShortVerDto
		for _, postReportResult := range reportResult.PostReports {
			postReportDtos = append(postReportDtos, ToPostReportShortVerDto(postReportResult))
		}
		return postReportDtos
	case consts.COMMENT_REPORT:
		var commentReportDtos []*CommentReportShortVerDto
		for _, commentReportResult := range reportResult.CommentReports {
			commentReportDtos = append(commentReportDtos, ToCommentReportShortVerDto(commentReportResult))
		}
		return commentReportDtos
	default:
		return nil
	}
}
