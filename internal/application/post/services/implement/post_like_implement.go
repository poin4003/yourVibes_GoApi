package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sPostLike struct {
	userRepo         postRepo.IUserRepository
	postRepo         postRepo.IPostRepository
	postLikeRepo     postRepo.ILikeUserPostRepository
	notificationRepo postRepo.INotificationRepository
}

func NewPostLikeImplement(
	userRepo postRepo.IUserRepository,
	postRepo postRepo.IPostRepository,
	postLikeRepo postRepo.ILikeUserPostRepository,
	notificationRepo postRepo.INotificationRepository,
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
	command *postCommand.LikePostCommand,
) (result *postCommand.LikePostCommandResult, err error) {
	result = &postCommand.LikePostCommandResult{}
	result.Post = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Find exist post
	postFound, err := s.postRepo.GetOne(ctx, command.PostId, command.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, fmt.Errorf("error when find post %v", err.Error())
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
		return result, fmt.Errorf("error when find user %v", err.Error())
	}

	// 3. Check like status (like or dislike)
	likeUserPostEntity, err := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)
	if err != nil {
		return result, err
	}

	checkLiked, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostEntity)
	if err != nil {
		return result, fmt.Errorf("failed to check like: %w", err)
	}

	// 4. Handle like and dislike
	if !checkLiked {
		// 4.1.1 Create new like if it not exist
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostEntity); err != nil {
			return result, fmt.Errorf("failed to create like: %w", err)
		}

		// 4.1.2. Plus 1 to likeCount of post
		updateData := &postEntity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount + 1),
		}

		err = updateData.ValidatePostUpdate()
		if err != nil {
			return result, fmt.Errorf("failed to create post: %w", err)
		}

		postUpdated, err := s.postRepo.UpdateOne(ctx, postFound.ID, updateData)

		// 4.1.3. Check like to response
		checkLikedToResponse, err := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.1.4. Push notification to owner of the post
		notification, err := notificationEntity.NewNotification(
			userLike.FamilyName+" "+userLike.Name,
			userLike.AvatarUrl,
			postFound.UserId,
			consts.LIKE_POST,
			(postFound.ID).String(),
			postFound.Content,
		)

		_, err = s.notificationRepo.CreateOne(ctx, notification)
		if err != nil {
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
			return result, fmt.Errorf("failed to send notification: %w", err)
		}

		// 4.1.6. Map to result
		result.Post = mapper.NewPostWithLikedParamResultFromEntity(postUpdated, isLiked)
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
			return result, fmt.Errorf("failed to delete like: %w", err)
		}

		// 4.2.2. Update -1 likeCount
		updateData := &postEntity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount - 1),
		}

		err = updateData.ValidatePostUpdate()
		if err != nil {
			return result, fmt.Errorf("failed to create post: %w", err)
		}

		postUpdated, err := s.postRepo.UpdateOne(ctx, postFound.ID, updateData)

		// 4.2.3. Check like to response
		checkLikedToResponse, err := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)
		if err != nil {
			return result, fmt.Errorf("failed to create post: %w", err)
		}

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.2.4. Map post to postDto
		result.Post = mapper.NewPostWithLikedParamResultFromEntity(postUpdated, isLiked)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK

		return result, nil
	}
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	query *postQuery.GetPostLikeQuery,
) (result *postQuery.GetPostLikeQueryResult, err error) {
	result = &postQuery.GetPostLikeQueryResult{}
	result.Users = nil
	result.PagingResponse = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError

	likeUserPostEntities, paging, err := s.postLikeRepo.GetLikeUserPost(ctx, query)
	if err != nil {
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
