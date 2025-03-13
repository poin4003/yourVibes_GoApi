package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewCommentResult(
	comment *reportEntity.CommentForReport,
) *common.CommentForReportResult {
	if comment == nil {
		return nil
	}

	return &common.CommentForReportResult{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		User:            NewUserResult(comment.User),
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		Status:          comment.Status,
	}
}
