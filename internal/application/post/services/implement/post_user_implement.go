package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	mapper2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/mapper"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
)

type sPostUser struct {
	userRepo         post_repo.IUserRepository
	FriendRepo       post_repo.IFriendRepository
	NewFeedRepo      post_repo.INewFeedRepository
	postRepo         post_repo.IPostRepository
	mediaRepo        post_repo.IMediaRepository
	likeUserPostRepo post_repo.ILikeUserPostRepository
	notificationRepo post_repo.INotificationRepository
}

func NewPostUserImplement(
	userRepo post_repo.IUserRepository,
	FriendRepo post_repo.IFriendRepository,
	NewFeedRepo post_repo.INewFeedRepository,
	postRepo post_repo.IPostRepository,
	mediaRepo post_repo.IMediaRepository,
	likeUserPostRepo post_repo.ILikeUserPostRepository,
	notificationRepo post_repo.INotificationRepository,
) *sPostUser {
	return &sPostUser{
		userRepo:         userRepo,
		FriendRepo:       FriendRepo,
		NewFeedRepo:      NewFeedRepo,
		postRepo:         postRepo,
		mediaRepo:        mediaRepo,
		likeUserPostRepo: likeUserPostRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	postModel *models.Post,
	inMedia []multipart.File,
) (post *models.Post, resultCode int, httpStatusCode int, err error) {
	// 1. CreatePost
	newPost, err := s.postRepo.CreatePost(ctx, postModel)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	// 2. Create Media and upload media to cloudinary_util
	if len(inMedia) > 0 {
		for _, file := range inMedia {
			// 2.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			if mediaUrl == "" {
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to upload media to cloudinary: empty media url")
			}

			// 2.2. create Media model and save to database
			mediaTemp := &models.Media{
				PostId:   newPost.ID,
				MediaUrl: mediaUrl,
			}

			_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
			if err != nil {
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	// 3. Find user
	userFound, err := s.userRepo.GetOne(ctx, "id=?", postModel.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find user: %w", err)
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to get user: %w", err)
	}

	// 4. Update post count for user
	userFound.PostCount++
	userUpdate := &user_entity.UserUpdate{
		PostCount: &userFound.PostCount,
	}
	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdate)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to update user: %w", err)
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to update user: %w", err)
	}

	// 5. Create new feed for user friend
	// 5.1. Get friend id of user friend list
	friendIds, err := s.FriendRepo.GetFriendIds(ctx, userFound.ID)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to get friends: %w", err)
	}

	// 5.2. If user don't have friend, return
	if len(friendIds) == 0 {
		return newPost, response.ErrCodeSuccess, http.StatusOK, nil
	}

	// 5.3. Create new feed for friend
	err = s.NewFeedRepo.CreateManyNewFeed(ctx, newPost.ID, friendIds)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create new feed: %w", err)
	}

	// 5.4. Create notification for friend
	//notificationEntities := make([]user_entity.Notification, len(friendIds))
	//for i, friendId := range friendIds {
	//	content := truncate.TruncateContent(newPost.Content, 20)
	//	notificationEntity, err := user_entity.NewNotification(
	//		userFound.FamilyName+" "+userFound.Name,
	//		userFound.AvatarUrl,
	//		friendId,
	//		consts.NEW_POST,
	//		newPost.ID.String(),
	//		content,
	//	)
	//
	//	if err != nil {
	//		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	//	}
	//
	//	notificationEntities[i] = *notificationEntity
	//}
	//
	//createdNotifications, err := s.notificationRepo.CreateMany(ctx, notificationEntities)
	//if err != nil {
	//	return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create notifications: %w", err)
	//}
	//
	//// 5.5. Send realtime notification (websocket)
	//notificationDto := mapper.MapNotificationToNotificationDto(createdNotifications[0])
	//userIds := make([]string, len(friendIds))
	//for i, friendId := range friendIds {
	//	userIds[i] = friendId.String()
	//}
	//
	//err = global.SocketHub.SendMultipleNotifications(userIds, notificationDto)
	//if err != nil {
	//	return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to send notifications: %w", err)
	//}

	return newPost, response.ErrCodeSuccess, http.StatusInternalServerError, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	postId uuid.UUID,
	updateData map[string]interface{},
	deleteMediaIds []uint,
	inMedia []multipart.File,
) (post *models.Post, resultCode int, httpStatusCode int, err error) {
	// 1. update post information
	postModel, err := s.postRepo.UpdatePost(ctx, postId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	// 2. delete media in database and delete media from cloudinary
	if len(deleteMediaIds) > 0 {
		for _, mediaId := range deleteMediaIds {
			// 2.1. Get media information from database
			media, err := s.mediaRepo.GetMedia(ctx, "id=?", mediaId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, response.ErrDataNotFound, http.StatusBadRequest, err
				}
				return nil, response.ErrDataNotFound, http.StatusInternalServerError, fmt.Errorf("failed to get media record: %w", err)
			}

			// 2.2. Delete media from cloudinary
			if media.MediaUrl != "" {
				if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
					return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete media record: %w", err)
				}
			}

			// 2.3. Delete media from databases
			if err := s.mediaRepo.DeleteMedia(ctx, mediaId); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, response.ErrDataNotFound, http.StatusBadRequest, nil
				}
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete media record: %w", err)
			}
		}
	}

	// 3. Create Media and upload media to cloudinary_util
	if len(inMedia) > 0 {
		for _, file := range inMedia {
			// 3.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			// 3.2. create Media model and save to database
			mediaTemp := &models.Media{
				PostId:   postId,
				MediaUrl: mediaUrl,
			}

			_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
			if err != nil {
				return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	return postModel, response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	postId uuid.UUID,
) (resultCode int, httpStatusCode int, err error) {
	// 1. Get media array of post
	medias, err := s.mediaRepo.GetManyMedia(ctx, "post_id=?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return response.ErrDataNotFound, http.StatusInternalServerError, fmt.Errorf("failed to get media records: %w", err)
	}

	// 2. Delete media from database and cloudinary
	for _, media := range medias {
		// 2.1. Delete media from cloudinary
		if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
			return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete media record: %w", err)
		}

		// 2.1. Delete media from databases
		if err := s.mediaRepo.DeleteMedia(ctx, media.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ErrDataNotFound, http.StatusBadRequest, nil
			}
			return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete media record: %w", err)
		}
	}

	postModel, err := s.postRepo.DeletePost(ctx, postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete media records: %w", err)
	}

	userFound, err := s.userRepo.GetOne(ctx, "id=?", postModel.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to get user: %w", err)
	}

	userFound.PostCount--
	userUpdateEntity := &user_entity.UserUpdate{PostCount: pointer.Ptr(userFound.PostCount)}
	newUserUpdateEntity, err := user_entity.NewUserUpdate(userUpdateEntity)
	if err != nil {
		return response.ErrCodeValidate, http.StatusBadRequest, fmt.Errorf("failed to validate: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, newUserUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return response.ErrDataNotFound, http.StatusInternalServerError, fmt.Errorf("failed to update media records: %w", err)
	}

	return response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	postId uuid.UUID,
	userId uuid.UUID,
) (postDto *response2.PostDto, resultCode int, httpStatusCode int, err error) {
	postModel, err := s.postRepo.GetPost(ctx, "id=?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, err
	}

	isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, &models.LikeUserPost{
		PostId: postId,
		UserId: userId,
	})

	postDto = mapper2.MapPostToPostDto(postModel, isLiked)

	return postDto, response.ErrCodeSuccess, http.StatusOK, nil
}

func (s *sPostUser) GetManyPosts(
	ctx context.Context,
	query *query.PostQueryObject,
	userId uuid.UUID,
) (postDtos []*response2.PostDto, resultCode int, httpStatusCode int, pagingResponse *response.PagingResponse, err error) {
	postModels, paging, err := s.postRepo.GetManyPost(ctx, query)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	for _, post := range postModels {
		isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, &models.LikeUserPost{
			PostId: post.ID,
			UserId: userId,
		})

		postDto := mapper2.MapPostToPostDto(post, isLiked)
		postDtos = append(postDtos, postDto)
	}

	return postDtos, response.ErrCodeSuccess, http.StatusOK, paging, nil
}
