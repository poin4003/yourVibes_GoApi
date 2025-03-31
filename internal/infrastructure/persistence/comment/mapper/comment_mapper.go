package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToCommentModel(comment *entities.Comment) *models.Comment {
	c := &models.Comment{
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		Status:          comment.Status,
	}
	c.ID = comment.ID

	return c
}

func ToUserEntity(
	user *models.User,
) *entities.User {
	if user == nil {
		return nil
	}

	return &entities.User{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}

func FromCommentModel(c *models.Comment) *entities.Comment {
	if c == nil {
		return nil
	}

	return &entities.Comment{
		ID:              c.ID,
		PostId:          c.PostId,
		UserId:          c.UserId,
		User:            ToUserEntity(&c.User),
		ParentId:        c.ParentId,
		Content:         c.Content,
		LikeCount:       c.LikeCount,
		RepCommentCount: c.RepCommentCount,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
		Status:          c.Status,
	}
}
