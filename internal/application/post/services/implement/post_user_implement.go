package implement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/poin4003/yourVibes_GoApi/pkg/utils/media"

	"github.com/poin4003/yourVibes_GoApi/global"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/truncate"
	"gorm.io/gorm"
)

type sPostUser struct {
	userRepo         postRepo.IUserRepository
	friendRepo       postRepo.IFriendRepository
	newFeedRepo      postRepo.INewFeedRepository
	postRepo         postRepo.IPostRepository
	mediaRepo        postRepo.IMediaRepository
	likeUserPostRepo postRepo.ILikeUserPostRepository
	notificationRepo postRepo.INotificationRepository
	advertiseRepo    postRepo.IAdvertiseRepository
}

func NewPostUserImplement(
	userRepo postRepo.IUserRepository,
	friendRepo postRepo.IFriendRepository,
	newFeedRepo postRepo.INewFeedRepository,
	postRepo postRepo.IPostRepository,
	mediaRepo postRepo.IMediaRepository,
	likeUserPostRepo postRepo.ILikeUserPostRepository,
	notificationRepo postRepo.INotificationRepository,
	advertiseRepo postRepo.IAdvertiseRepository,
) *sPostUser {
	return &sPostUser{
		userRepo:         userRepo,
		friendRepo:       friendRepo,
		newFeedRepo:      newFeedRepo,
		postRepo:         postRepo,
		mediaRepo:        mediaRepo,
		likeUserPostRepo: likeUserPostRepo,
		notificationRepo: notificationRepo,
		advertiseRepo:    advertiseRepo,
	}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	command *postCommand.CreatePostCommand,
) (result *postCommand.CreatePostCommandResult, err error) {
	result = &postCommand.CreatePostCommandResult{
		Post:           nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. CreatePost
	newPost, err := postEntity.NewPost(
		command.UserId,
		command.Content,
		command.Privacy,
		command.Location,
	)

	if err != nil {
		return result, err
	}

	postCreated, err := s.postRepo.CreateOne(ctx, newPost)
	if err != nil {
		return result, err
	}

	// 2. Create Media and upload media
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 2.1. Save file and get mediaUrl
			mediaUrl, err := media.SaveMedia(&file)

			if err != nil {
				return result, fmt.Errorf("failed to upload media: %w", err)
			}

			// 2.2. create Media model and save to database
			mediaEntity, err := postEntity.NewMedia(postCreated.ID, mediaUrl)
			if err != nil {
				return result, err
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				return result, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	// 3. Find user
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("failed to find user: %w", err)
		}
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	// 4. Update post count for user
	userFound.PostCount++
	userUpdate := &userEntity.UserUpdate{
		PostCount: &userFound.PostCount,
	}
	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("failed to update user: %w", err)
		}
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	// 5. Check privacy of post
	if postCreated.Privacy == consts.PRIVATE {
		result.Post = mapper.NewPostResultFromEntity(postCreated)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 6. Create new feed for user friend
	// 6.1. Get friend id of user friend list
	friendIds, err := s.friendRepo.GetFriendIds(ctx, userFound.ID)
	if err != nil {
		return result, fmt.Errorf("failed to get friends: %w", err)
	}

	// 6.2. If user don't have friend, return
	if len(friendIds) == 0 {
		result.Post = mapper.NewPostResultFromEntity(postCreated)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}

	// 6.3. Create new feed for friend
	err = s.newFeedRepo.CreateMany(ctx, newPost.ID, userFound.ID)
	if err != nil {
		return result, fmt.Errorf("failed to create new feed: %w", err)
	}

	// 6.4. Create notification for friend
	var notificationEntities []*notificationEntity.Notification
	for _, friendId := range friendIds {
		content := truncate.TruncateContent(newPost.Content, 20)
		notification, err := notificationEntity.NewNotification(
			userFound.FamilyName+" "+userFound.Name,
			userFound.AvatarUrl,
			friendId,
			consts.NEW_POST,
			newPost.ID.String(),
			content,
		)

		if err != nil {
			return result, err
		}

		notificationEntities = append(notificationEntities, notification)
	}

	_, err = s.notificationRepo.CreateMany(ctx, notificationEntities)
	if err != nil {
		return result, fmt.Errorf("failed to create notifications: %w", err)
	}

	// 6.5. Send realtime notification (websocket)
	for _, friendId := range friendIds {
		notificationSocketResponse := &consts.NotificationSocketResponse{
			From:             userFound.FamilyName + " " + userFound.Name,
			FromUrl:          userFound.AvatarUrl,
			UserId:           friendId,
			NotificationType: consts.NEW_POST,
			ContentId:        (postCreated.ID).String(),
		}

		err = global.SocketHub.SendNotification(friendId.String(), notificationSocketResponse)
		if err != nil {
			return result, fmt.Errorf("failed to send notifications: %w", err)
		}
	}

	// 7. Validate post after create
	validatePost, err := postValidator.NewValidatedPost(postCreated)
	if err != nil {
		return result, fmt.Errorf("failed to validate post: %w", err)
	}

	result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	command *postCommand.UpdatePostCommand,
) (result *postCommand.UpdatePostCommandResult, err error) {
	result = &postCommand.UpdatePostCommandResult{
		Post:           nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. update post information
	updateData := &postEntity.PostUpdate{
		Content:  command.Content,
		Privacy:  command.Privacy,
		Location: command.Location,
	}

	err = updateData.ValidatePostUpdate()
	if err != nil {
		return result, err
	}

	postUpdated, err := s.postRepo.UpdateOne(ctx, *command.PostId, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. delete media in database and delete media
	if len(command.MediaIDs) > 0 {
		for _, mediaId := range command.MediaIDs {
			// 2.1. Get media information from database
			mediaRecord, err := s.mediaRepo.GetOne(ctx, "id=?", mediaId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					result.ResultCode = response.ErrDataNotFound
					result.HttpStatusCode = http.StatusBadRequest
					return result, err
				}
				return result, fmt.Errorf("failed to get media record: %w", err)
			}

			// 2.2. Delete media from cloudinary
			if mediaRecord.MediaUrl != "" {
				if err := media.DeleteMedia(mediaRecord.MediaUrl); err != nil {
					return result, fmt.Errorf("failed to delete media record: %w", err)
				}
			}

			// 2.3. Delete media from databases
			if err := s.mediaRepo.DeleteOne(ctx, mediaId); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					result.ResultCode = response.ErrDataNotFound
					result.HttpStatusCode = http.StatusBadRequest
					return result, nil
				}
				return result, fmt.Errorf("failed to delete media record: %w", err)
			}
		}
	}

	// 3. Create Media and upload media
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 3.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := media.SaveMedia(&file)

			if err != nil {
				return result, fmt.Errorf("failed to upload media: %w", err)
			}

			// 3.2. create Media model and save to database
			mediaEntity, err := postEntity.NewMedia(postUpdated.ID, mediaUrl)
			if err != nil {
				return result, err
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				return result, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	result.Post = mapper.NewPostResultFromEntity(postUpdated)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	command *postCommand.DeletePostCommand,
) (result *postCommand.DeletePostCommandResult, err error) {
	result = &postCommand.DeletePostCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get media array of post
	medias, err := s.mediaRepo.GetMany(ctx, "post_id=?", command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("failed to get media records: %w", err)
	}

	// 2. Delete media from database and folder
	for _, mediaRecord := range medias {
		// 2.1. Delete media from folder
		if err := media.DeleteMedia(mediaRecord.MediaUrl); err != nil {
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to delete media record: %w", err)
		}

		// 2.1. Delete media from databases
		if err := s.mediaRepo.DeleteOne(ctx, mediaRecord.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, nil
			}
			return result, fmt.Errorf("failed to delete media record: %w", err)
		}
	}

	deleteCondition := map[string]interface{}{
		"post_id": command.PostId,
	}

	// 3. Delete new feed
	err = s.newFeedRepo.DeleteMany(ctx, deleteCondition)
	if err != nil {
		return result, fmt.Errorf("failed to delete media records: %w", err)
	}

	// 4. Delete advertise and bill
	err = s.advertiseRepo.DeleteMany(ctx, deleteCondition)
	if err != nil {
		return result, fmt.Errorf("failed to delete media records: %w", err)
	}

	// 5. Delete post
	post, err := s.postRepo.DeleteOne(ctx, *command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("failed to delete media records: %w", err)
	}

	// 6. Find user
	userFound, err := s.userRepo.GetOne(ctx, "id=?", post.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	// 7. Update post count of user
	userFound.PostCount--

	userUpdateEntity := &userEntity.UserUpdate{PostCount: pointer.Ptr(userFound.PostCount)}

	err = userUpdateEntity.ValidateUserUpdate()
	if err != nil {
		return result, fmt.Errorf("failed to update user: %w", err)
	}

	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("failed to update media records: %w", err)
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	query *postQuery.GetOnePostQuery,
) (result *postQuery.GetOnePostQueryResult, err error) {
	result = &postQuery.GetOnePostQueryResult{
		Post:           nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get post
	postFound, err := s.postRepo.GetOne(ctx, query.PostId, query.AuthenticatedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check privacy
	isOwner := postFound.UserId == query.AuthenticatedUserId
	if !isOwner {
		switch postFound.Privacy {
		case consts.PUBLIC:
		case consts.FRIEND_ONLY:
			isFriend, err := s.friendRepo.CheckFriendExist(ctx, &userEntity.Friend{
				UserId:   postFound.UserId,
				FriendId: query.AuthenticatedUserId,
			})
			if err != nil {
				return result, err
			}
			if !isFriend {
				result.Post = nil
				result.ResultCode = response.ErrPostFriendAccess
				result.HttpStatusCode = http.StatusBadRequest
				return result, fmt.Errorf("authenticated user is not a friend")
			}
		case consts.PRIVATE:
			result.Post = nil
			result.ResultCode = response.ErrPostPrivateAccess
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("only owner can access this post")
		default:
			result.Post = nil
			result.ResultCode = response.ErrPostPrivateAccess
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("only owner can access this post")
		}
	}

	// 4. Return
	result.Post = mapper.NewPostWithLikedResultFromEntity(postFound)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostUser) GetManyPosts(
	ctx context.Context,
	query *postQuery.GetManyPostQuery,
) (result *postQuery.GetManyPostQueryResult, err error) {
	result = &postQuery.GetManyPostQueryResult{
		Posts:          nil,
		PagingResponse: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	postEntities, paging, err := s.postRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var postResults []*common.PostResultWithLiked
	for _, post := range postEntities {
		postResult := mapper.NewPostWithLikedResultFromEntity(post)
		postResults = append(postResults, postResult)
	}

	result.Posts = postResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sPostUser) CheckPostOwner(
	ctx context.Context,
	query *postQuery.CheckPostOwnerQuery,
) (result *postQuery.CheckPostOwnerQueryResult, err error) {
	result = &postQuery.CheckPostOwnerQueryResult{
		IsOwner:        false,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}

	isOwner, err := s.postRepo.CheckPostOwner(ctx, query.PostId, query.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	result.IsOwner = isOwner
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
