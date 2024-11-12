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

func FromPostModel(postModel *models.Post) *post_entity.Post {
	var user = &post_entity.User{
		ID:         postModel.User.ID,
		FamilyName: postModel.User.FamilyName,
		Name:       postModel.User.Name,
		AvatarUrl:  postModel.User.AvatarUrl,
	}

	var parentPost = &post_entity.Post{}
	if postModel.ParentPost != nil {
		parentPost = &post_entity.Post{
			ID:              postModel.ParentPost.ID,
			UserId:          postModel.ParentPost.UserId,
			User:            user,
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
		}
	} else {
		parentPost = &post_entity.Post{}
	}

	var medias = []*post_entity.Media{}
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

	var post = &post_entity.Post{
		ID:              postModel.ID,
		UserId:          postModel.UserId,
		User:            user,
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
	return post
}
