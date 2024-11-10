package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/mapper"
)

func MapPostToPostDto(post *models.Post, isLiked bool) *response.PostDto {
	var parentPost *response.ParentPostDto

	if post.ParentPost != nil {
		parentPost = &response.ParentPostDto{
			ID:              post.ParentPost.ID,
			UserId:          post.ParentPost.UserId,
			User:            mapper.MapUserToUserDtoShortVer(&post.ParentPost.User),
			Content:         post.ParentPost.Content,
			LikeCount:       post.ParentPost.LikeCount,
			CommentCount:    post.ParentPost.CommentCount,
			Privacy:         post.ParentPost.Privacy,
			Location:        post.ParentPost.Location,
			IsAdvertisement: post.ParentPost.IsAdvertisement,
			Status:          post.ParentPost.Status,
			IsLiked:         isLiked,
			CreatedAt:       post.ParentPost.CreatedAt,
			UpdatedAt:       post.ParentPost.UpdatedAt,
			DeletedAt:       post.ParentPost.DeletedAt,
			Media:           post.ParentPost.Media,
		}
	}

	return &response.PostDto{
		ID:              post.ID,
		UserId:          post.UserId,
		User:            mapper.MapUserToUserDtoShortVer(&post.User),
		ParentId:        post.ParentId,
		ParentPost:      parentPost,
		Content:         post.Content,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		Privacy:         post.Privacy,
		Location:        post.Location,
		IsAdvertisement: post.IsAdvertisement,
		Status:          post.Status,
		IsLiked:         isLiked,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		DeletedAt:       post.DeletedAt,
		Media:           post.Media,
	}
}

func MapPostToUpdatedPostDto(post *models.Post) *response.UpdatedPostDto {
	return &response.UpdatedPostDto{
		ID:              post.ID,
		UserId:          post.UserId,
		User:            mapper.MapUserToUserDtoShortVer(&post.User),
		ParentId:        post.ParentId,
		ParentPost:      post.ParentPost,
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

func MapPostToNewPostDto(post *models.Post) *response.NewPostDto {
	return &response.NewPostDto{
		ID:              post.ID,
		UserId:          post.UserId,
		ParentId:        post.ParentId,
		ParentPost:      post.ParentPost,
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
	}
}

func MapToPostFromCreateDto(
	input *request.CreatePostInput,
	userId uuid.UUID,
) *models.Post {
	return &models.Post{
		UserId:   userId,
		Content:  input.Content,
		Privacy:  input.Privacy,
		Location: input.Location,
	}
}

func MapToPostFromUpdateDto(
	input *request.UpdatePostInput,
) map[string]interface{} {
	updateData := make(map[string]interface{})

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
