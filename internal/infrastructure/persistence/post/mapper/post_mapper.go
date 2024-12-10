package mapper

import (
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToPostModel(post *post_entity.Post) *models.Post {
	p := &models.Post{
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
	}

	return p
}

func ToUserEntity(
	user *models.User,
) *post_entity.User {
	if user == nil {
		return nil
	}

	return &post_entity.User{
		ID:         user.ID,
		FamilyName: user.FamilyName,
		Name:       user.Name,
		AvatarUrl:  user.AvatarUrl,
	}
}

func FromPostModel(postModel *models.Post) *post_entity.Post {
	var parentPost *post_entity.Post
	if postModel.ParentPost != nil {
		var medias []*post_entity.Media
		for _, media := range postModel.ParentPost.Media {
			medias = append(medias, &post_entity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &post_entity.Post{
			ID:              postModel.ParentPost.ID,
			UserId:          postModel.ParentPost.UserId,
			User:            ToUserEntity(&postModel.ParentPost.User),
			ParentId:        postModel.ParentPost.ParentId,
			Content:         postModel.ParentPost.Content,
			LikeCount:       postModel.ParentPost.LikeCount,
			CommentCount:    postModel.ParentPost.CommentCount,
			Privacy:         postModel.ParentPost.Privacy,
			Location:        postModel.ParentPost.Location,
			IsAdvertisement: postModel.ParentPost.IsAdvertisement,
			Status:          postModel.ParentPost.Status,
			CreatedAt:       postModel.ParentPost.CreatedAt,
			UpdatedAt:       postModel.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*post_entity.Media
	for _, media := range postModel.Media {
		medias = append(medias, &post_entity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	return &post_entity.Post{
		ID:              postModel.ID,
		UserId:          postModel.UserId,
		User:            ToUserEntity(&postModel.User),
		ParentId:        postModel.ParentId,
		ParentPost:      parentPost,
		Content:         postModel.Content,
		LikeCount:       postModel.LikeCount,
		CommentCount:    postModel.CommentCount,
		Privacy:         postModel.Privacy,
		Location:        postModel.Location,
		IsAdvertisement: postModel.IsAdvertisement,
		Status:          postModel.Status,
		CreatedAt:       postModel.CreatedAt,
		UpdatedAt:       postModel.UpdatedAt,
		Media:           medias,
	}
}

func FromPostWithLikedModel(postModel *models.PostWithLiked) *post_entity.PostWithLiked {
	var parentPost *post_entity.Post
	if postModel.ParentPost != nil {
		var medias []*post_entity.Media
		for _, media := range postModel.ParentPost.Media {
			medias = append(medias, &post_entity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &post_entity.Post{
			ID:              postModel.ParentPost.ID,
			UserId:          postModel.ParentPost.UserId,
			User:            ToUserEntity(&postModel.ParentPost.User),
			ParentId:        postModel.ParentPost.ParentId,
			Content:         postModel.ParentPost.Content,
			LikeCount:       postModel.ParentPost.LikeCount,
			CommentCount:    postModel.ParentPost.CommentCount,
			Privacy:         postModel.ParentPost.Privacy,
			Location:        postModel.ParentPost.Location,
			IsAdvertisement: postModel.ParentPost.IsAdvertisement,
			Status:          postModel.ParentPost.Status,
			CreatedAt:       postModel.ParentPost.CreatedAt,
			UpdatedAt:       postModel.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*post_entity.Media
	for _, media := range postModel.Media {
		medias = append(medias, &post_entity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	return &post_entity.PostWithLiked{
		ID:              postModel.ID,
		UserId:          postModel.UserId,
		User:            ToUserEntity(&postModel.User),
		ParentId:        postModel.ParentId,
		ParentPost:      parentPost,
		Content:         postModel.Content,
		LikeCount:       postModel.LikeCount,
		CommentCount:    postModel.CommentCount,
		Privacy:         postModel.Privacy,
		Location:        postModel.Location,
		IsAdvertisement: postModel.IsAdvertisement,
		IsLiked:         postModel.IsLiked,
		Status:          postModel.Status,
		CreatedAt:       postModel.CreatedAt,
		UpdatedAt:       postModel.UpdatedAt,
		Media:           medias,
	}
	return nil
}
