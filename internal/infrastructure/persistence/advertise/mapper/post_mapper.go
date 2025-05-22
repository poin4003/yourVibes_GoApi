package mapper

import (
	advetiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func FromPostModel(postModel *models.Post) *advetiseEntity.PostForAdvertise {
	if postModel == nil {
		return nil
	}

	var parentPost *advetiseEntity.PostForAdvertise
	if postModel.ParentPost != nil {
		var medias []*advetiseEntity.Media
		for _, media := range postModel.ParentPost.Media {
			medias = append(medias, &advetiseEntity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &advetiseEntity.PostForAdvertise{
			ID:              postModel.ParentPost.ID,
			UserId:          postModel.ParentPost.UserId,
			User:            FromUserModel(&postModel.ParentPost.User),
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

	var medias []*advetiseEntity.Media
	for _, media := range postModel.Media {
		medias = append(medias, &advetiseEntity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	return &advetiseEntity.PostForAdvertise{
		ID:              postModel.ID,
		UserId:          postModel.UserId,
		User:            FromUserModel(&postModel.User),
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

func FromPostModelToShortPostEntity(postModel *models.Post) *advetiseEntity.ShortPostForAdvertise {
	if postModel == nil {
		return nil
	}

	var parentPost *advetiseEntity.ShortPostForAdvertise
	if postModel.ParentPost != nil {
		var medias []*advetiseEntity.Media
		for _, media := range postModel.ParentPost.Media {
			medias = append(medias, &advetiseEntity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &advetiseEntity.ShortPostForAdvertise{
			ID:      postModel.ParentPost.ID,
			Content: postModel.ParentPost.Content,
			Media:   medias,
		}
	}

	var medias []*advetiseEntity.Media
	for _, media := range postModel.Media {
		medias = append(medias, &advetiseEntity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	return &advetiseEntity.ShortPostForAdvertise{
		ID:         postModel.ID,
		ParentPost: parentPost,
		Content:    postModel.Content,
		Media:      medias,
	}
}
