package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	comment_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/validator"
)

func NewCommentResultFromValidateEntity(
	comment *comment_validator.ValidatedComment,
) *common.CommentResult {
	return NewCommentResultFromEntity(&comment.Comment)
}

func NewCommentResultFromEntity(
	comment *entities.Comment,
) *common.CommentResult {
	if comment == nil {
		return nil
	}

	return &common.CommentResult{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		User:            NewUserResultFromEntity(comment.User),
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		Status:          comment.Status,
	}
}

func NewCommentWithLikedResultFromEntity(
	comment *entities.Comment,
) *common.CommentResultWithLiked {
	if comment == nil {
		return nil
	}

	return &common.CommentResultWithLiked{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		User:            NewUserResultFromEntity(comment.User),
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		Status:          comment.Status,
		IsLiked:         comment.IsLiked,
	}
}

func NewCommentWithLikedResultFromEntityAndIsLiked(
	comment *entities.Comment,
	isLiked bool,
) *common.CommentResultWithLiked {
	if comment == nil {
		return nil
	}

	return &common.CommentResultWithLiked{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		User:            NewUserResultFromEntity(comment.User),
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		Status:          comment.Status,
		IsLiked:         isLiked,
	}
}
