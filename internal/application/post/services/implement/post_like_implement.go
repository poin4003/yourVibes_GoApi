package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repository"
	entities2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/entities"
	mapper2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/post/post_user/dto/mapper"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/rest/user/user_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
	likeUserPostModel *entities2.LikeUserPost,
	userId uuid.UUID,
) (postDto *response2.PostDto, resultCode int, httpStatusCode int, err error) {
	// 1. Find exist post
	postFound, err := s.postRepo.GetPost(ctx, "id=?", likeUserPostModel.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find post %w", err.Error())
	}

	// 2. Find exist user
	userLike, err := s.userRepo.GetUser(ctx, "id=?", userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrDataNotFound, http.StatusBadRequest, err
		}
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("Error when find user %w", err.Error())
	}

	// 3. Check like status (like or dislike)
	checkLiked, err := s.postLikeRepo.CheckUserLikePost(ctx, likeUserPostModel)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to check like: %w", err)
	}

	// 4. Handle like and dislike
	if !checkLiked {
		// 4.1.1 Create new like if it not exist
		if err := s.postLikeRepo.CreateLikeUserPost(ctx, likeUserPostModel); err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create like: %w", err)
		}

		// 4.1.2. Plus 1 to likeCount of post
		postFound.LikeCount++
		_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
			"like_count": postFound.LikeCount,
		})

		// 4.1.3. Check if Authenticated User liked the post
		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, &entities2.LikeUserPost{
			PostId: postFound.ID,
			UserId: userId,
		})

		// 4.1.4. Push notification to owner of the post
		notificationModel := &entities2.Notification{
			From:             userLike.FamilyName + " " + userLike.Name,
			FromUrl:          userLike.AvatarUrl,
			UserId:           postFound.UserId,
			NotificationType: consts.LIKE_POST,
			ContentId:        (postFound.ID).String(),
			Content:          postFound.Content,
		}
		notification, err := s.notificationRepo.CreateNotification(ctx, notificationModel)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to create notification: %w", err)
		}

		// 4.1.5. Send realtime notification (websocket)
		notificationDto := mapper.MapNotificationToNotificationDto(notification)

		err = global.SocketHub.SendNotification(postFound.UserId.String(), notificationDto)
		if err != nil {
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to send notification: %w", err)
		}

		// 4.1.6. Map Post to PostDto to response for client
		postDto = mapper2.MapPostToPostDto(postFound, isLiked)

		// 4.1.7. Response for controller
		return postDto, response.ErrCodeSuccess, http.StatusOK, nil
	} else {
		// 4.2.1. Delete like if it exits
		if err := s.postLikeRepo.DeleteLikeUserPost(ctx, likeUserPostModel); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.ErrDataNotFound, http.StatusBadRequest, fmt.Errorf("failed to find delete like: %w", err)
			}
			return nil, response.ErrServerFailed, http.StatusInternalServerError, fmt.Errorf("failed to delete like: %w", err)
		}

		// 4.2.2. Update -1 likeCount
		postFound.LikeCount--
		_, err = s.postRepo.UpdatePost(ctx, postFound.ID, map[string]interface{}{
			"like_count": postFound.LikeCount,
		})

		// 4.2.3. Check if Authenticated User liked the post
		isLiked, _ := s.postLikeRepo.CheckUserLikePost(ctx, &entities2.LikeUserPost{
			PostId: postFound.ID,
			UserId: userId,
		})

		// 4.2.4. Map post to postDto
		postDto = mapper2.MapPostToPostDto(postFound, isLiked)

		// 4.2.5. Response for controller
		return postDto, response.ErrCodeSuccess, http.StatusOK, nil
	}
}

func (s *sPostLike) GetUsersOnLikes(
	ctx context.Context,
	postId uuid.UUID,
	query *query.PostLikeQueryObject,
) (users []*entities2.User, resultCode int, httpStatusCode int, responsePaging *response.PagingResponse, err error) {
	likeUserPostModel, paging, err := s.postLikeRepo.GetLikeUserPost(ctx, postId, query)
	if err != nil {
		return nil, response.ErrServerFailed, http.StatusInternalServerError, nil, err
	}

	return likeUserPostModel, response.ErrCodeSuccess, http.StatusOK, paging, nil
}
