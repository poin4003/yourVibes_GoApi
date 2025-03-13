package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"

	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func FromPostModel(postModel *models.Post) *reportEntity.PostForReport {
	if postModel == nil {
		return nil
	}

	var parentPost *reportEntity.PostForReport
	if postModel.ParentPost != nil {
		var medias []*reportEntity.Media
		for _, media := range postModel.ParentPost.Media {
			medias = append(medias, &reportEntity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &reportEntity.PostForReport{
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

	var medias []*reportEntity.Media
	for _, media := range postModel.Media {
		medias = append(medias, &reportEntity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	return &reportEntity.PostForReport{
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
