package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
)

type CommentForReportDto struct {
	ID              uuid.UUID         `json:"id"`
	PostId          uuid.UUID         `json:"post_id"`
	UserId          uuid.UUID         `json:"user_id"`
	User            *UserForReportDto `json:"user"`
	ParentId        *uuid.UUID        `json:"parent_id"`
	Content         string            `json:"content"`
	LikeCount       int               `json:"like_count"`
	RepCommentCount int               `json:"rep_comment_count"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Status          bool              `json:"status"`
}

func ToCommentForReportDto(commentResult *common.CommentForReportResult) *CommentForReportDto {
	return &CommentForReportDto{
		ID:              commentResult.ID,
		PostId:          commentResult.PostId,
		UserId:          commentResult.UserId,
		User:            ToUserForReportDto(commentResult.User),
		ParentId:        commentResult.ParentId,
		Content:         commentResult.Content,
		LikeCount:       commentResult.LikeCount,
		RepCommentCount: commentResult.RepCommentCount,
		CreatedAt:       commentResult.CreatedAt,
		UpdatedAt:       commentResult.UpdatedAt,
		Status:          commentResult.Status,
	}
}
