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

func MapToCommentFromUpdateDto(
	input *comment_dto.UpdateCommentInput,
) map[string]interface{} {
	updateData := make(map[string]interface{})

	if input.Content != nil {
		updateData["content"] = *input.Content
	}

	return updateData
}

func MapCommentToCommentDto(
	comment *model.Comment,
	isLiked bool,
) *comment_dto.CommentDto {
	return &comment_dto.CommentDto{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		IsLiked:         isLiked,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		User:            MapUserToUserDtoShortVer(&comment.User),
	}
}

func MapCommentToNewCommentDto(
	comment *model.Comment,
) *comment_dto.NewCommentDto {
	return &comment_dto.NewCommentDto{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
	}
}

func MapCommentToUpdatedCommentDto(
	comment *model.Comment,
) *comment_dto.UpdatedCommentDto {
	return &comment_dto.UpdatedCommentDto{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		User:            MapUserToUserDtoShortVer(&comment.User),
	}
}
