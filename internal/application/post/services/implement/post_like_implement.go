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
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sPostLike struct {
	userRepo         post_repo.IUserRepository
	postRepo         post_repo.IPostRepository
	postLikeRepo     post_repo.ILikeUserPostRepository
	notificationRepo post_repo.INotificationRepository
}

func NewPostLikeImplement(
	userRepo post_repo.IUserRepository,
	postRepo post_repo.IPostRepository,
	postLikeRepo post_repo.ILikeUserPostRepository,
	notificationRepo post_repo.INotificationRepository,
) *sPostLike {
	return &sPostLike{
		userRepo:         userRepo,
		postRepo:         postRepo,
		postLikeRepo:     postLikeRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sPostLike) LikePost(
	ctx context.Context,
	command *post_command.LikePostCommand,
) (result *post_command.LikePostCommandResult, err error) {
	result = &post_command.LikePostCommandResult{}
	// 1. Find exist post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
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
		return result, fmt.Errorf("Error when find post %w", err.Error())
	}

	// 2. Find exist user
	userLike, err := s.userRepo.GetById(ctx, command.UserId)
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
		return result, fmt.Errorf("Error when find user %w", err.Error())
	}

	// 3. Check like status (like or dislike)
	likeUserPostEntity, err := post_entity.NewLikeUserPostEntity(command.UserId, command.PostId)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	checkLiked, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostEntity)
	if err != nil {
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, fmt.Errorf("failed to check like: %w", err)
	}

	// 4. Handle like and dislike
	if !checkLiked {
		// 4.1.1 Create new like if it not exist
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostEntity); err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to create like: %w", err)
		}

		// 4.1.2. Plus 1 to likeCount of post
		updateData := &post_entity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount + 1),
		}

		updatePost, err := post_entity.NewPostUpdate(updateData)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to create post: %w", err)
		}

		_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)

		// 4.1.3. Check like to response
		checkLikedToResponse, err := post_entity.NewLikeUserPostEntity(command.UserId, command.PostId)

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.1.4. Push notification to owner of the post
		notificationEntity, err := notification_entity.NewNotification(
			userLike.FamilyName+" "+userLike.Name,
			userLike.AvatarUrl,
			postFound.UserId,
			consts.LIKE_POST,
			(postFound.ID).String(),
			postFound.Content,
		)

		_, err = s.notificationRepo.CreateOne(ctx, notificationEntity)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to create notification: %w", err)
		}

		// 4.1.5. Send realtime notification (websocket)
		notificationSocketResponse := &consts.NotificationSocketResponse{
			From:             userLike.FamilyName + " " + userLike.Name,
			FromUrl:          userLike.AvatarUrl,
			UserId:           postFound.UserId,
			NotificationType: consts.LIKE_POST,
			ContentId:        (postFound.ID).String(),
		}

		err = global.SocketHub.SendNotification(postFound.UserId.String(), notificationSocketResponse)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to send notification: %w", err)
		}

		// 4.1.6. Map Post to PostDto to response for client
		result.Post = mapper.NewPostWithLikedResultFromEntity(postFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK

		// 4.1.7. Response for controller
		return result, nil
	} else {
		// 4.2.1. Delete like if it exits
		if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostEntity); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Post = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, fmt.Errorf("failed to find delete like: %w", err)
			}
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to delete like: %w", err)
		}

		// 4.2.2. Update -1 likeCount
		updateData := &post_entity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount - 1),
		}

		updatePost, err := post_entity.NewPostUpdate(updateData)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, fmt.Errorf("failed to create post: %w", err)
		}

		_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)

		// 4.2.3. Check like to response
		checkLikedToResponse, err := post_entity.NewLikeUserPostEntity(command.UserId, command.PostId)

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.2.4. Map post to postDto
		result.Post = mapper.NewPostWithLikedResultFromEntity(postFound, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK

		return result, nil
	}
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	query *post_query.GetPostLikeQuery,
) (result *post_query.GetPostLikeQueryResult, err error) {
	result = &post_query.GetPostLikeQueryResult{}

	likeUserPostEntities, paging, err := s.postLikeRepo.GetLikeUserPost(ctx, query)
	if err != nil {
		result.Users = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
		return result, err
	}

	var likeUserPostResults []*common.UserResult
	for _, likeUserPostEntity := range likeUserPostEntities {
		likeUserPostResults = append(likeUserPostResults, mapper.NewUserResultFromEntity(likeUserPostEntity))
	}

	result.Users = likeUserPostResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}
