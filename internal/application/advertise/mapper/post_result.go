package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	advertiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

func NewPostResult(
	post *advertiseEntity.PostForAdvertise,
) *common.PostForAdvertiseResult {
	if post == nil {
		return nil
	}

	var parentPost *common.PostForAdvertiseResult
	if post.ParentPost != nil {
		parentPost = &common.PostForAdvertiseResult{
			ID:              post.ParentPost.ID,
			UserId:          post.ParentPost.UserId,
			User:            NewUserResult(post.ParentPost.User),
			ParentId:        post.ParentPost.ParentId,
			Content:         post.ParentPost.Content,
			LikeCount:       post.ParentPost.LikeCount,
			CommentCount:    post.ParentPost.CommentCount,
			Privacy:         post.ParentPost.Privacy,
			Location:        post.ParentPost.Location,
			IsAdvertisement: post.ParentPost.IsAdvertisement,
			Status:          post.ParentPost.Status,
			CreatedAt:       post.ParentPost.CreatedAt,
			UpdatedAt:       post.ParentPost.UpdatedAt,
			Media:           NewMediaResult(post.ParentPost.Media),
		}
	}

	return &common.PostForAdvertiseResult{
		ID:              post.ID,
		UserId:          post.UserId,
		User:            NewUserResult(post.User),
		ParentId:        post.ParentId,
		ParentPost:      parentPost,
		Content:         post.Content,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		Privacy:         post.Privacy,
		Location:        post.Location,
		IsAdvertisement: post.IsAdvertisement,
		Status:          post.Status,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		Media:           NewMediaResult(post.Media),
	}
}
