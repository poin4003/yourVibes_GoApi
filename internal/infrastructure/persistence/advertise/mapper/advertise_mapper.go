package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToAdvertiseModel(advertise *entities.Advertise) *models.Advertise {
	a := &models.Advertise{
		PostId:    advertise.PostId,
		StartDate: advertise.StartDate,
		EndDate:   advertise.EndDate,
		CreatedAt: advertise.CreatedAt,
		UpdatedAt: advertise.UpdatedAt,
	}
	a.ID = advertise.ID

	return a
}

func FromAdvertiseModel(a *models.Advertise) *entities.Advertise {
	var post = &entities.PostForAdvertise{
		User: ToUserForAdvertiseEntity(&a.Post.User),
	}

	var bill = &entities.Bill{
		ID:          a.Bill.ID,
		AdvertiseId: a.Bill.AdvertiseId,
		Price:       a.Bill.Price,
		CreatedAt:   a.Bill.CreatedAt,
		UpdateAt:    a.Bill.UpdatedAt,
		Status:      a.Bill.Status,
	}

	var advertise = &entities.Advertise{
		PostId:    a.PostId,
		Post:      post,
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Bill:      bill,
	}
	advertise.ID = a.ID

	return advertise
}

func FromAdvertiseModelForAdvertiseDetail(a *models.Advertise) *entities.Advertise {
	// var parentPost *entities.PostForAdvertise
	// if a.Post.ParentPost != nil {
	// 	var medias []*entities.Media
	// 	for _, media := range a.Post.ParentPost.Media {
	// 		medias = append(medias, &entities.Media{
	// 			ID:        media.ID,
	// 			MediaUrl:  media.MediaUrl,
	// 			PostId:    media.PostId,
	// 			Status:    media.Status,
	// 			CreatedAt: media.CreatedAt,
	// 			UpdatedAt: media.UpdatedAt,
	// 		})
	// 	}
	// 	parentPost = &entities.PostForAdvertise{
	// 		ID:              a.Post.ParentPost.ID,
	// 		UserId:          a.Post.ParentPost.UserId,
	// 		User:            ToUserForAdvertiseEntity(&a.Post.ParentPost.User),
	// 		ParentId:        a.Post.ParentPost.ParentId,
	// 		Content:         a.Post.ParentPost.Content,
	// 		LikeCount:       a.Post.ParentPost.LikeCount,
	// 		CommentCount:    a.Post.ParentPost.CommentCount,
	// 		Privacy:         a.Post.ParentPost.Privacy,
	// 		Location:        a.Post.ParentPost.Location,
	// 		IsAdvertisement: a.Post.ParentPost.IsAdvertisement,
	// 		Status:          a.Post.ParentPost.Status,
	// 		CreatedAt:       a.Post.ParentPost.CreatedAt,
	// 		UpdatedAt:       a.Post.ParentPost.UpdatedAt,
	// 		Media:           medias,
	// 	}
	// }

	// var medias []*entities.Media
	// for _, media := range a.Post.Media {
	// 	medias = append(medias, &entities.Media{
	// 		ID:        media.ID,
	// 		MediaUrl:  media.MediaUrl,
	// 		PostId:    media.PostId,
	// 		Status:    media.Status,
	// 		CreatedAt: media.CreatedAt,
	// 		UpdatedAt: media.UpdatedAt,
	// 	})
	// }

	// post := &entities.PostForAdvertise{
	// 	ID:              a.Post.ID,
	// 	UserId:          a.Post.UserId,
	// 	User:            ToUserForAdvertiseEntity(&a.Post.User),
	// 	ParentId:        a.Post.ParentId,
	// 	ParentPost:      parentPost,
	// 	Content:         a.Post.Content,
	// 	LikeCount:       a.Post.LikeCount,
	// 	CommentCount:    a.Post.CommentCount,
	// 	Privacy:         a.Post.Privacy,
	// 	Location:        a.Post.Location,
	// 	IsAdvertisement: a.Post.IsAdvertisement,
	// 	Status:          a.Post.Status,
	// 	CreatedAt:       a.Post.CreatedAt,
	// 	UpdatedAt:       a.Post.UpdatedAt,
	// 	Media:           medias,
	// }

	// var bill = &entities.Bill{
	// 	ID:          a.Bill.ID,
	// 	AdvertiseId: a.Bill.AdvertiseId,
	// 	Price:       a.Bill.Price,
	// 	CreatedAt:   a.Bill.CreatedAt,
	// 	UpdateAt:    a.Bill.UpdatedAt,
	// 	Status:      a.Bill.Status,
	// }

	// var advertise = &entities.Advertise{
	// 	PostId:    a.PostId,
	// 	Post:      post,
	// 	StartDate: a.StartDate,
	// 	EndDate:   a.EndDate,
	// 	CreatedAt: a.CreatedAt,
	// 	UpdatedAt: a.UpdatedAt,
	// 	Bill:      bill,
	// }
	// advertise.ID = a.ID

	// return advertiseS

	return nil
}

func ToUserForAdvertiseEntity(
	userModel *models.User,
) *entities.UserForAdvertise {
	if userModel == nil {
		return nil
	}

	var userForAdvertise = &entities.UserForAdvertise{
		FamilyName:   userModel.FamilyName,
		Name:         userModel.Name,
		Email:        userModel.Email,
		Password:     userModel.Password,
		PhoneNumber:  userModel.PhoneNumber,
		Birthday:     userModel.Birthday,
		AvatarUrl:    userModel.AvatarUrl,
		CapwallUrl:   userModel.CapwallUrl,
		Privacy:      userModel.Privacy,
		Biography:    userModel.Biography,
		AuthType:     userModel.AuthType,
		AuthGoogleId: userModel.AuthGoogleId,
		PostCount:    userModel.PostCount,
		FriendCount:  userModel.FriendCount,
		Status:       userModel.Status,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
	userForAdvertise.ID = userModel.ID

	return userForAdvertise
}
