package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/truncate"
	"go.uber.org/zap"
)

type sPostShare struct {
	userRepo           repository.IUserRepository
	postRepo           repository.IPostRepository
	mediaRepo          repository.IMediaRepository
	newFeedRepo        repository.INewFeedRepository
	friendRepo         repository.IFriendRepository
	postCache          cache.IPostCache
	postEventPublisher *producer.PostEventPublisher
}

func NewPostShareImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	mediaRepo repository.IMediaRepository,
	newFeedRepo repository.INewFeedRepository,
	friendRepo repository.IFriendRepository,
	postCache cache.IPostCache,
	postEventPublisher *producer.PostEventPublisher,
) *sPostShare {
	return &sPostShare{
		userRepo:           userRepo,
		postRepo:           postRepo,
		mediaRepo:          mediaRepo,
		newFeedRepo:        newFeedRepo,
		friendRepo:         friendRepo,
		postCache:          postCache,
		postEventPublisher: postEventPublisher,
	}
}

func (s *sPostShare) SharePost(
	ctx context.Context,
	command *postCommand.SharePostCommand,
) (result *postCommand.SharePostCommandResult, err error) {
	result = &postCommand.SharePostCommandResult{
		Post: nil,
	}

	// 1. Find post by post_id
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return nil, response.NewDataNotFoundError("post not found")
	}

	rootPostId := command.PostId
	if postFound.ParentId != nil {
		rootPost, err := s.postRepo.GetById(ctx, *postFound.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
		if rootPost == nil {
			return nil, response.NewDataNotFoundError("post not found")
		}
		rootPostId = rootPost.ID
	}

	// 3. Copy post info from root post
	newPost, err := postEntity.NewPostForShare(
		command.UserId,
		command.Content,
		command.Privacy,
		command.Location,
		&rootPostId,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}
	// 4. Create new post
	newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	validatePost, err := postValidator.NewValidatedPost(newSharePost)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 5. Get friend list
	friendIds, err := s.friendRepo.GetFriendIds(ctx, validatePost.UserId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 6. Return if user don't have friend
	if len(friendIds) == 0 {
		result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
		return result, nil
	}

	// 7. Create new feed for friend
	err = s.newFeedRepo.CreateMany(ctx, newSharePost.ID, newSharePost.User.ID)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 8. Create notification for friend
	notification, err := notificationEntity.NewNotification(
		newSharePost.User.FamilyName+" "+newSharePost.User.Name,
		newSharePost.User.AvatarUrl,
		newSharePost.User.ID,
		consts.NEW_SHARE,
		newSharePost.ID.String(),
		truncate.TruncateContent(newSharePost.Content, 20),
	)
	if err != nil {
		global.Logger.Error("Failed to create notification entity", zap.Error(err))
		return result, nil
	}

	// 9. Publish to RabbitMQ to handle Notification
	notiMsg := mapper.NewNotificationResult(notification)
	if err = s.postEventPublisher.PublishNotification(ctx, notiMsg, "notification.bulk.db_websocket"); err != nil {
		global.Logger.Error("Failed to publish notification for friend", zap.Error(err))
	}

	// 10. Delete cache
	s.postCache.DeleteFeeds(ctx, consts.RK_USER_FEED, validatePost.UserId)
	s.postCache.DeleteFeeds(ctx, consts.RK_PERSONAL_POST, validatePost.UserId)
	s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)

	return &postCommand.SharePostCommandResult{
		Post: mapper.NewPostResultFromValidateEntity(validatePost),
	}, nil
}
