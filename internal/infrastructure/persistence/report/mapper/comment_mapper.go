package mapper

import (
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromCommentModel(c *models.Comment) *reportEntity.CommentForReport {
	if c == nil {
		return nil
	}

	return &reportEntity.CommentForReport{
		ID:              c.ID,
		PostId:          c.PostId,
		UserId:          c.UserId,
		User:            FromUserModel(&c.User),
		ParentId:        c.ParentId,
		Content:         c.Content,
		LikeCount:       c.LikeCount,
		RepCommentCount: c.RepCommentCount,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
		Status:          c.Status,
	}
}
