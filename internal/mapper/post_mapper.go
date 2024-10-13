package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/post_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

func MapPostToPostDto(post *model.Post) *post_dto.PostDto {
	return &post_dto.PostDto{
		ID:              post.ID,
		UserId:          post.UserId,
		User:            post.User,
		ParentId:        post.ParentId,
		ParentPost:      post.ParentPost,
		Title:           post.Title,
		Content:         post.Content,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		Privacy:         post.Privacy,
		Location:        post.Location,
		IsAdvertisement: post.IsAdvertisement,
		Status:          post.Status,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		DeletedAt:       post.DeletedAt,
		Media:           post.Media,
	}
}

func MapToPostFromCreateDto(
	input *post_dto.CreatePostInput,
	userId uuid.UUID,
) *model.Post {
	return &model.Post{
		UserId:   userId,
		Title:    input.Title,
		Content:  input.Content,
		Privacy:  input.Privacy,
		Location: input.Location,
	}
}

func MapToPostFromUpdateDto(
	input *post_dto.UpdatePostInput,
) map[string]interface{} {
	updateData := make(map[string]interface{})

	if input.Title != nil {
		updateData["title"] = *input.Title
	}
	if input.Content != nil {
		updateData["content"] = *input.Content
	}
	if input.Privacy != nil {
		updateData["privacy"] = *input.Privacy
	}
	if input.Location != nil {
		updateData["location"] = *input.Location
	}

	return updateData
}
