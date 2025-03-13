package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
)

type CommentReportDto struct {
	ReportId          uuid.UUID            `json:"report_id"`
	UserId            uuid.UUID            `json:"user_id"`
	ReportedCommentId uuid.UUID            `json:"reported_comment_id"`
	AdminId           *uuid.UUID           `json:"admin_id"`
	User              *UserForReportDto    `json:"user"`
	ReportedComment   *CommentForReportDto `json:"reported_comment"`
	Post              *PostForReportDto    `json:"post"`
	Admin             *AdminDto            `json:"admin"`
	Reason            string               `json:"reason"`
	Status            bool                 `json:"status"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

type CommentReportShortVerDto struct {
	ReportId          uuid.UUID  `json:"report_id"`
	UserId            uuid.UUID  `json:"user_id"`
	ReportedCommentId uuid.UUID  `json:"reported_comment_id"`
	UserEmail         string     `json:"user_email"`
	AdminEmail        *string    `json:"admin_email"`
	AdminId           *uuid.UUID `json:"admin_id"`
	Reason            string     `json:"reason"`
	Status            bool       `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func ToCommentReportDto(
	commentReportResult *common.CommentReportResult,
) *CommentReportDto {
	if commentReportResult == nil {
		return nil
	}

	return &CommentReportDto{
		ReportId:          commentReportResult.ReportId,
		UserId:            commentReportResult.UserId,
		ReportedCommentId: commentReportResult.ReportedCommentId,
		AdminId:           commentReportResult.AdminId,
		User:              ToUserForReportDto(commentReportResult.User),
		ReportedComment:   ToCommentForReportDto(commentReportResult.ReportedComment),
		Post:              ToPostForReportDto(commentReportResult.Post),
		Admin:             ToAdminDto(commentReportResult.Admin),
		Reason:            commentReportResult.Reason,
		Status:            commentReportResult.Status,
		CreatedAt:         commentReportResult.CreatedAt,
		UpdatedAt:         commentReportResult.UpdatedAt,
	}
}

func ToCommentReportShortVerDto(
	commentReportResult *common.CommentReportShortVerResult,
) *CommentReportShortVerDto {
	if commentReportResult == nil {
		return nil
	}

	return &CommentReportShortVerDto{
		ReportId:          commentReportResult.ReportId,
		UserId:            commentReportResult.UserId,
		ReportedCommentId: commentReportResult.ReportedCommentId,
		UserEmail:         commentReportResult.UserEmail,
		AdminEmail:        commentReportResult.AdminEmail,
		AdminId:           commentReportResult.AdminId,
		Reason:            commentReportResult.Reason,
		Status:            commentReportResult.Status,
		CreatedAt:         commentReportResult.CreatedAt,
		UpdatedAt:         commentReportResult.UpdatedAt,
	}
}
