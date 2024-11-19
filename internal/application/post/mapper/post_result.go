package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

func NewPostWithLikedResultFromEntity(
	post *entities.Post,
	isLiked bool,
) *common.PostResultWithLiked {
	if post == nil {
		return nil
	}

	var parentPost *common.PostResultWithLiked

	if post.ParentPost != nil {
		parentPost = &common.PostResultWithLiked{
			ID:              post.ParentPost.ID,
			UserId:          post.ParentPost.UserId,
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
			User:            NewUserResultFromEntity(post.User),
		}
	}

	return &common.PostResultWithLiked{
		ID:              post.ID,
		UserId:          post.UserId,
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
		Media:           NewMediaResultsFromEntity(post.Media),
		User:            NewUserResultFromEntity(post.User),
	}
}

func NewPostResultFromEntity(
	post *entities.Post,
) *common.PostResult {
	if post == nil {
		return nil
	}

	return &common.PostResult{
		ID:              post.ID,
		UserId:          post.UserId,
		ParentId:        post.ParentId,
		Content:         post.Content,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		Privacy:         post.Privacy,
		Location:        post.Location,
		IsAdvertisement: post.IsAdvertisement,
		Status:          post.Status,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		Media:           NewMediaResultsFromEntity(post.Media),
		User:            NewUserResultFromEntity(post.User),
	}
}
