package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"time"
)

type CommentReportDto struct {
	UserId            uuid.UUID
	ReportedCommentId uuid.UUID
	User              *UserForReportDto
	ReportedComment   *CommentForReportDto
	Post              *PostForReportDto
	Reason            string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func ToCommentReportDto(commentResult *common.CommentReportResult) *CommentReportDto {
	return &CommentReportDto{
		UserId:            commentResult.UserId,
		ReportedCommentId: commentResult.ReportedCommentId,
		User:              ToUserForReportDto(commentResult.User),
		ReportedComment:   ToCommentForReportDto(commentResult.ReportedComment),
		Post:              ToPostForReportDto(commentResult.Post),
		Reason:            commentResult.Reason,
		Status:            commentResult.Status,
		CreatedAt:         commentResult.CreatedAt,
		UpdatedAt:         commentResult.UpdatedAt,
	}
}
