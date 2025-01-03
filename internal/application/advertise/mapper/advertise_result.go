package mapper

import (
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertiseValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/validator"
)

func NewAdvertiseWithBillResultFromValidateEntity(
	advertise *advertiseValidator.ValidateAdvertise,
) *common.AdvertiseWithBillResult {
	return NewAdvertiseWithBillResultFromEntity(&advertise.Advertise)
}

func NewAdvertiseWithBillResultFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseWithBillResult {
	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	return &common.AdvertiseWithBillResult{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		UserEmail:    advertise.Post.User.Email,
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
		Bill:         NewBillWithoutAdvertiseResultFromEntity(advertise.Bill),
	}
}

func NewAdvertiseWithoutBillResultFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseWithoutBillResult {
	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	return &common.AdvertiseWithoutBillResult{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
	}
}

func NewAdvertiseDetailFromEntity(
	advertise *entities.Advertise,
) *common.AdvertiseDetail {

	if advertise == nil {
		return nil
	}

	duration := advertise.EndDate.Sub(time.Now())
	dayRemaining := int(duration.Hours() / 24)

	if dayRemaining == 0 && duration.Hours() > 0 {
		dayRemaining = 1
	}

	if dayRemaining < 0 {
		dayRemaining = 0
	}

	var parentPost *common.PostForAdvertiseResult
	if advertise.Post.ParentPost != nil {
		var medias []*common.MediaResult
		for _, media := range advertise.Post.ParentPost.Media {
			medias = append(medias, &common.MediaResult{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &common.PostForAdvertiseResult{
			ID:              advertise.Post.ParentPost.ID,
			UserId:          advertise.Post.ParentPost.UserId,
			User:            NewUserForAdvertiseResult(advertise.Post.ParentPost.User),
			ParentId:        advertise.Post.ParentPost.ParentId,
			Content:         advertise.Post.ParentPost.Content,
			LikeCount:       advertise.Post.ParentPost.LikeCount,
			CommentCount:    advertise.Post.ParentPost.CommentCount,
			Privacy:         advertise.Post.ParentPost.Privacy,
			Location:        advertise.Post.ParentPost.Location,
			IsAdvertisement: advertise.Post.ParentPost.IsAdvertisement,
			Status:          advertise.Post.ParentPost.Status,
			CreatedAt:       advertise.Post.ParentPost.CreatedAt,
			UpdatedAt:       advertise.Post.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*common.MediaResult
	for _, media := range advertise.Post.Media {
		medias = append(medias, &common.MediaResult{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	post := &common.PostForAdvertiseResult{
		ID:              advertise.Post.ID,
		UserId:          advertise.Post.UserId,
		User:            NewUserForAdvertiseResult(advertise.Post.User),
		ParentId:        advertise.Post.ParentId,
		ParentPost:      parentPost,
		Content:         advertise.Post.Content,
		LikeCount:       advertise.Post.LikeCount,
		CommentCount:    advertise.Post.CommentCount,
		Privacy:         advertise.Post.Privacy,
		Location:        advertise.Post.Location,
		IsAdvertisement: advertise.Post.IsAdvertisement,
		Status:          advertise.Post.Status,
		CreatedAt:       advertise.Post.CreatedAt,
		UpdatedAt:       advertise.Post.UpdatedAt,
		Media:           medias,
	}

	return &common.AdvertiseDetail{
		ID:           advertise.ID,
		PostId:       advertise.PostId,
		Post:         post,
		StartDate:    advertise.StartDate,
		EndDate:      advertise.EndDate,
		DayRemaining: dayRemaining,
		CreatedAt:    advertise.CreatedAt,
		UpdatedAt:    advertise.UpdatedAt,
		Bill:         NewBillWithoutAdvertiseResultFromEntity(advertise.Bill),
	}
}
