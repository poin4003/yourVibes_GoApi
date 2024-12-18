package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"time"
)

type CommentReportDto struct {
	UserId            uuid.UUID            `json:"user_id"`
	ReportedCommentId uuid.UUID            `json:"reported_comment_id"`
	User              *UserForReportDto    `json:"user"`
	ReportedComment   *CommentForReportDto `json:"reported_comment"`
	Post              *PostForReportDto    `json:"post"`
	Reason            string               `json:"reason"`
	Status            bool                 `json:"status"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
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
