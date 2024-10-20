package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/comment_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapToCommentFromCreateDto(
	input *comment_dto.CreateCommentInput,
	userId uuid.UUID,
) *model.Comment {
	return &model.Comment{
		UserId:   userId,
		PostId:   input.PostId,
		ParentId: input.ParentId,
		Content:  input.Content,
	}
}

func MapCommentToCommentDto(
	comment *model.Comment,
) *comment_dto.CommentDto {
	return &comment_dto.CommentDto{
		ID:        comment.ID,
		PostId:    comment.PostId,
		UserId:    comment.UserId,
		ParentId:  comment.ParentId,
		Content:   comment.Content,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User:      MapUserToUserDtoShortVer(&comment.User),
	}
}
