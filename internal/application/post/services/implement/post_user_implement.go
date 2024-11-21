package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notification_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	post_validator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/truncate"
	"gorm.io/gorm"
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
	command *post_command.CreatePostCommand,
) (result *post_command.CreatePostCommandResult, err error) {
	result = &post_command.CreatePostCommandResult{}
	// 1. CreatePost
	newPost, err := post_entity.NewPost(
		command.UserId,
		command.Content,
		command.Privacy,
		command.Location,
	)

	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	postEntity, err := s.postRepo.CreateOne(ctx, newPost)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Create Media and upload media to cloudinary_util
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 2.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			if mediaUrl == "" {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to upload media to cloudinary: empty media url")
			}

			// 2.2. create Media model and save to database
			mediaEntity, err := post_entity.NewMedia(postEntity.ID, mediaUrl)
			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, err
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	// 3. Find user
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("failed to find user: %w", err)
		}
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	// 4. Update post count for user
	userFound.PostCount++
	userUpdate := &user_entity.UserUpdate{
		PostCount: &userFound.PostCount,
	}
	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("failed to update user: %w", err)
		}
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	// 5. Create new feed for user friend
	// 5.1. Get friend id of user friend list
	friendIds, err := s.FriendRepo.GetFriendIds(ctx, userFound.ID)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to get friends: %w", err)
	}

	// 5.2. If user don't have friend, return
	if len(friendIds) == 0 {
		result.Post = mapper.NewPostResultFromEntity(postEntity)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 5.3. Create new feed for friend
	err = s.NewFeedRepo.CreateMany(ctx, newPost.ID, friendIds)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to create new feed: %w", err)
	}

	// 5.4. Create notification for friend
	var notificationEntities []*notification_entity.Notification
	for _, friendId := range friendIds {
		content := truncate.TruncateContent(newPost.Content, 20)
		notificationEntity, err := notification_entity.NewNotification(
			userFound.FamilyName+" "+userFound.Name,
			userFound.AvatarUrl,
			friendId,
			consts.NEW_POST,
			newPost.ID.String(),
			content,
		)

		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}

		notificationEntities = append(notificationEntities, notificationEntity)
	}

	_, err = s.notificationRepo.CreateMany(ctx, notificationEntities)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to create notifications: %w", err)
	}

	// 5.5. Send realtime notification (websocket)
	for _, friendId := range friendIds {
		notificationSocketResponse := &consts.NotificationSocketResponse{
			From:             userFound.FamilyName + " " + userFound.Name,
			FromUrl:          userFound.AvatarUrl,
			UserId:           friendId,
			NotificationType: consts.NEW_POST,
			ContentId:        (postEntity.ID).String(),
		}

		err = global.SocketHub.SendNotification(friendId.String(), notificationSocketResponse)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to send notifications: %w", err)
		}
	}

	// 6. Validate post after create
	validatePost, err := post_validator.NewValidatedPost(postEntity)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to validate post: %w", err)
	}

	result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	command *post_command.UpdatePostCommand,
) (result *post_command.UpdatePostCommandResult, err error) {
	result = &post_command.UpdatePostCommandResult{}
	// 1. update post information
	updateData := &post_entity.PostUpdate{
		Content:  command.Content,
		Privacy:  command.Privacy,
		Location: command.Location,
	}

	err = updateData.ValidatePostUpdate()
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	postEntity, err := s.postRepo.UpdateOne(ctx, *command.PostId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. delete media in database and delete media from cloudinary
	if len(command.MediaIDs) > 0 {
		for _, mediaId := range command.MediaIDs {
			// 2.1. Get media information from database
			media, err := s.mediaRepo.GetOne(ctx, "id=?", mediaId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					result.Post = nil
					result.ResultCode = response.ErrDataNotFound
					result.HttpStatusCode = http.StatusBadRequest
					return result, err
				}
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to get media record: %w", err)
			}

			// 2.2. Delete media from cloudinary
			if media.MediaUrl != "" {
				if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
					result.Post = nil
					result.ResultCode = response.ErrServerFailed
					result.HttpStatusCode = http.StatusInternalServerError
					return result, fmt.Errorf("failed to delete media record: %w", err)
				}
			}

			// 2.3. Delete media from databases
			if err := s.mediaRepo.DeleteOne(ctx, mediaId); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					result.Post = nil
					result.ResultCode = response.ErrDataNotFound
					result.HttpStatusCode = http.StatusBadRequest
					return result, nil
				}
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to delete media record: %w", err)
			}
		}
	}

	// 3. Create Media and upload media to cloudinary_util
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 3.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			// 3.2. create Media model and save to database
			mediaEntity, err := post_entity.NewMedia(postEntity.ID, mediaUrl)
			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, err
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				result.Post = nil
				result.ResultCode = response.ErrServerFailed
				result.HttpStatusCode = http.StatusInternalServerError
				return result, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	result.Post = mapper.NewPostResultFromEntity(postEntity)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	command *post_command.DeletePostCommand,
) (result *post_command.DeletePostCommandResult, err error) {
	result = &post_command.DeletePostCommandResult{}
	// 1. Get media array of post
	medias, err := s.mediaRepo.GetMany(ctx, "post_id=?", command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to get media records: %w", err)
	}

	// 2. Delete media from database and cloudinary
	for _, media := range medias {
		// 2.1. Delete media from cloudinary
		if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to delete media record: %w", err)
		}

		// 2.1. Delete media from databases
		if err := s.mediaRepo.DeleteOne(ctx, media.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, nil
			}
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to delete media record: %w", err)
		}
	}

	// 3. Delete post
	postEntity, err := s.postRepo.DeleteOne(ctx, *command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to delete media records: %w", err)
	}

	// 5. Find user
	userFound, err := s.userRepo.GetOne(ctx, "id=?", postEntity.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	// 6. Update post count of user
	userFound.PostCount--

	userUpdateEntity := &user_entity.UserUpdate{PostCount: pointer.Ptr(userFound.PostCount)}

	err = userUpdateEntity.ValidateUserUpdate()
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to update media records: %w", err)
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	query *post_query.GetOnePostQuery,
) (result *post_query.GetOnePostQueryResult, err error) {
	result = &post_query.GetOnePostQueryResult{}

	// 1. Get post
	postEntity, err := s.postRepo.GetOne(ctx, "id=?", query.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Check isLiked by authenticated user
	likeUserPostEntity, err := post_entity.NewLikeUserPostEntity(query.AuthenticatedUserId, query.PostId)

	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, likeUserPostEntity)

	// 3. Return
	result.Post = mapper.NewPostWithLikedResultFromEntity(postEntity, isLiked)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) GetManyPosts(
	ctx context.Context,
	query *post_query.GetManyPostQuery,
) (result *post_query.GetManyPostQueryResult, err error) {
	result = &post_query.GetManyPostQueryResult{}

	postEntities, paging, err := s.postRepo.GetMany(ctx, query)
	if err != nil {
		result.Posts = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
		return result, err
	}

	var postResults []*common.PostResultWithLiked
	for _, postEntity := range postEntities {
		likeUserPost, err := post_entity.NewLikeUserPostEntity(query.AuthenticatedUserId, postEntity.ID)
		if err != nil {
			result.Posts = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			result.PagingResponse = nil
			return result, err
		}

		isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, likeUserPost)

		postResult := mapper.NewPostWithLikedResultFromEntity(postEntity, isLiked)
		postResults = append(postResults, postResult)
	}

	result.Posts = postResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}
