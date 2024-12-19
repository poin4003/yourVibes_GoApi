package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"time"
)

type CommentReportDto struct {
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
	UserId            uuid.UUID  `json:"user_id"`
	ReportedCommentId uuid.UUID  `json:"reported_comment_id"`
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
	userReportResult *common.CommentReportShortVerResult,
) *CommentReportShortVerDto {
	if userReportResult == nil {
		return nil
	}

	return &CommentReportShortVerDto{
		UserId:            userReportResult.UserId,
		ReportedCommentId: userReportResult.ReportedCommentId,
		AdminId:           userReportResult.AdminId,
		Reason:            userReportResult.Reason,
		Status:            userReportResult.Status,
		CreatedAt:         userReportResult.CreatedAt,
		UpdatedAt:         userReportResult.UpdatedAt,
	}
}
