package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"

	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"go.uber.org/zap"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sPostLike struct {
	userRepo           repository.IUserRepository
	postRepo           repository.IPostRepository
	postLikeRepo       repository.ILikeUserPostRepository
	postCache          cache.IPostCache
	postEventPublisher *producer.PostEventPublisher
}

func NewPostLikeImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	postLikeRepo repository.ILikeUserPostRepository,
	postCache cache.IPostCache,
	postEventPublisher *producer.PostEventPublisher,
) *sPostLike {
	return &sPostLike{
		userRepo:           userRepo,
		postRepo:           postRepo,
		postLikeRepo:       postLikeRepo,
		postCache:          postCache,
		postEventPublisher: postEventPublisher,
	}
}

func (s *sPostLike) LikePost(
	ctx context.Context,
	command *postCommand.LikePostCommand,
) (result *postCommand.LikePostCommandResult, err error) {
	result = &postCommand.LikePostCommandResult{
		Post: nil,
	}
	// 1. Find exist post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return nil, response.NewDataNotFoundError("post not found")
	}

	// 2. Find exist user
	userLike, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
	}

	if userLike == nil {
		return nil, response.NewDataNotFoundError("user like this post not found")
	}

	// 3. Check like status (like or dislike)
	likeUserPostEntity, err := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	checkLiked, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	s.postCache.DeletePost(ctx, command.PostId)

	// 4. Handle like and dislike
	if !checkLiked {
		// 4.1.1 Create new like if it not exist
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostEntity); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 4.1.2. Plus 1 to likeCount of post
		updateData := &postEntity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount + 1),
		}

		err = updateData.ValidatePostUpdate()
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		postUpdated, _ := s.postRepo.UpdateOne(ctx, postFound.ID, updateData)

		// 4.1.3. Check like to response
		checkLikedToResponse, _ := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.1.4. Return if owner like his posts
		if userLike.ID == postFound.UserId {
			result.Post = mapper.NewPostWithLikedResultFromEntity(postUpdated, isLiked)
			return result, nil
		}

		// 4.1.5. Publish to RabbitMQ handle Notification
		notification, _ := notificationEntity.NewNotification(
			userLike.FamilyName+" "+userLike.Name,
			userLike.AvatarUrl,
			postFound.UserId,
			consts.LIKE_POST,
			(postFound.ID).String(),
			postFound.Content,
		)

		notifMsg := mapper.NewNotificationResult(notification)
		if err = s.postEventPublisher.PublishNotification(ctx, notifMsg, "notification.single.db_websocket"); err != nil {
			global.Logger.Error("Failed to publish notification result", zap.Error(err))
		}

		// 4.1.7. Map to result
		result.Post = mapper.NewPostWithLikedResultFromEntity(postUpdated, isLiked)
		// 4.1.8. Response for controller
		return result, nil
	} else {
		// 4.2.1. Delete like if it exits
		if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostEntity); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 4.2.2. Update -1 likeCount
		updateData := &postEntity.PostUpdate{
			LikeCount: pointer.Ptr(postFound.LikeCount - 1),
		}

		err = updateData.ValidatePostUpdate()
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		postUpdated, _ := s.postRepo.UpdateOne(ctx, postFound.ID, updateData)

		// 4.2.3. Check like to response
		checkLikedToResponse, err := postEntity.NewLikeUserPostEntity(command.UserId, command.PostId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, checkLikedToResponse)

		// 4.2.4. Map post to postDto
		result.Post = mapper.NewPostWithLikedResultFromEntity(postUpdated, isLiked)

		return result, nil
	}
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	query *postQuery.GetPostLikeQuery,
) (result *postQuery.GetPostLikeQueryResult, err error) {
	likeUserPostEntities, paging, err := s.postLikeRepo.GetLikeUserPost(ctx, query)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	var likeUserPostResults []*common.UserResult
	for _, likeUserPostEntity := range likeUserPostEntities {
		likeUserPostResults = append(likeUserPostResults, mapper.NewUserResultFromEntity(likeUserPostEntity))
	}

	return &postQuery.GetPostLikeQueryResult{
		Users:          likeUserPostResults,
		PagingResponse: paging,
	}, nil
}
