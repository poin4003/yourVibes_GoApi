package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/comment/comment_user/dto/response"
	mapper2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/mapper"
)

func MapToCommentFromCreateDto(
	input *request.CreateCommentInput,
	userId uuid.UUID,
) *models.Comment {
	return &models.Comment{
		UserId:   userId,
		PostId:   input.PostId,
		ParentId: input.ParentId,
		Content:  input.Content,
	}
}

func MapToCommentFromUpdateDto(
	input *request.UpdateCommentInput,
) map[string]interface{} {
	updateData := make(map[string]interface{})

	if input.Content != nil {
		updateData["content"] = *input.Content
	}

	return updateData
}

func MapCommentToCommentDto(
	comment *models.Comment,
	isLiked bool,
) *response.CommentDto {
	return &response.CommentDto{
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
		User:            mapper2.MapUserToUserDtoShortVer(&comment.User),
	}
}

func MapCommentToNewCommentDto(
	comment *models.Comment,
) *response.NewCommentDto {
	return &response.NewCommentDto{
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
	comment *models.Comment,
) *response.UpdatedCommentDto {
	return &response.UpdatedCommentDto{
		ID:              comment.ID,
		PostId:          comment.PostId,
		UserId:          comment.UserId,
		ParentId:        comment.ParentId,
		Content:         comment.Content,
		LikeCount:       comment.LikeCount,
		RepCommentCount: comment.RepCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		User:            mapper2.MapUserToUserDtoShortVer(&comment.User),
	}
}
