package implement

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	userEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/media"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/truncate"
	"go.uber.org/zap"
)

type sPostUser struct {
	userRepo           repositories.IUserRepository
	friendRepo         repositories.IFriendRepository
	newFeedRepo        repositories.INewFeedRepository
	postRepo           repositories.IPostRepository
	mediaRepo          repositories.IMediaRepository
	likeUserPostRepo   repositories.ILikeUserPostRepository
	advertiseRepo      repositories.IAdvertiseRepository
	postCache          cache.IPostCache
	postEventPublisher *producer.PostEventPublisher
}

func NewPostUserImplement(
	userRepo repositories.IUserRepository,
	friendRepo repositories.IFriendRepository,
	newFeedRepo repositories.INewFeedRepository,
	postRepo repositories.IPostRepository,
	mediaRepo repositories.IMediaRepository,
	likeUserPostRepo repositories.ILikeUserPostRepository,
	advertiseRepo repositories.IAdvertiseRepository,
	postCache cache.IPostCache,
	postEventPublisher *producer.PostEventPublisher,
) *sPostUser {
	return &sPostUser{
		userRepo:           userRepo,
		friendRepo:         friendRepo,
		newFeedRepo:        newFeedRepo,
		postRepo:           postRepo,
		mediaRepo:          mediaRepo,
		likeUserPostRepo:   likeUserPostRepo,
		advertiseRepo:      advertiseRepo,
		postCache:          postCache,
		postEventPublisher: postEventPublisher,
	}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	command *postCommand.CreatePostCommand,
) (result *postCommand.CreatePostCommandResult, err error) {
	result = &postCommand.CreatePostCommandResult{
		Post: nil,
	}
	// 1. CreatePost
	newPost, err := postEntity.NewPost(
		command.UserId,
		command.Content,
		command.Privacy,
		command.Location,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	postCreated, err := s.postRepo.CreateOne(ctx, newPost)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Create Media and upload media
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 2.1. Save file and get mediaUrl
			mediaUrl, err := media.SaveMedia(&file)

			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}

			// 2.2. create Media model and save to database
			mediaEntity, err := postEntity.NewMedia(postCreated.ID, mediaUrl)
			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}
		}
	}

	// 3. Find user
	userFound, err := s.userRepo.GetById(ctx, command.UserId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if userFound == nil {
		return nil, response.NewDataNotFoundError("user not found")
	}

	// 4. Update post count for user
	userFound.PostCount++
	userUpdate := &userEntity.UserUpdate{
		PostCount: &userFound.PostCount,
	}
	_, err = s.userRepo.UpdateOne(ctx, userFound.ID, userUpdate)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 5. Check privacy of post
	if postCreated.Privacy == consts.PRIVATE {
		result.Post = mapper.NewPostResultFromEntity(postCreated)
		return result, nil
	}

	// 6. Create new feed for user friend
	// 6.1. Get friend id of user friend list
	friendIds, err := s.friendRepo.GetFriendIds(ctx, userFound.ID)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 6.2. If user don't have friend, return
	if len(friendIds) == 0 {
		result.Post = mapper.NewPostResultFromEntity(postCreated)
		return result, nil
	}

	// 7. Create new feed for friend
	err = s.newFeedRepo.CreateMany(ctx, newPost.ID, userFound.ID)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 8. Create notification for friend
	notification, err := notificationEntity.NewNotification(
		userFound.FamilyName+" "+userFound.Name,
		userFound.AvatarUrl,
		userFound.ID,
		consts.NEW_POST,
		newPost.ID.String(),
		truncate.TruncateContent(newPost.Content, 20),
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

	// 10. Validate post after create
	validatePost, err := postValidator.NewValidatedPost(postCreated)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 11. Delete feed cache
	s.postCache.DeleteFeeds(ctx, consts.RK_USER_FEED, userFound.ID)
	s.postCache.DeleteFeeds(ctx, consts.RK_PERSONAL_POST, userFound.ID)
	s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friendIds)

	result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
	return result, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	command *postCommand.UpdatePostCommand,
) (result *postCommand.UpdatePostCommandResult, err error) {
	postFound, err := s.postRepo.GetById(ctx, *command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return nil, response.NewDataNotFoundError("post not found")
	}

	// 1. update post information
	updateData := &postEntity.PostUpdate{
		Content:  command.Content,
		Privacy:  command.Privacy,
		Location: command.Location,
	}

	err = updateData.ValidatePostUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	postUpdated, err := s.postRepo.UpdateOne(ctx, *command.PostId, updateData)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. delete media in database and delete media
	if len(command.MediaIDs) > 0 {
		for _, mediaId := range command.MediaIDs {
			// 2.1. Get media information from database
			mediaRecord, err := s.mediaRepo.GetOne(ctx, "id=?", mediaId)
			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}

			if mediaRecord == nil {
				return nil, response.NewDataNotFoundError("media not found")
			}

			// 2.2. Delete media from cloudinary
			if mediaRecord.MediaUrl != "" {
				if err := media.DeleteMedia(mediaRecord.MediaUrl); err != nil {
					return nil, response.NewServerFailedError(err.Error())
				}
			}

			// 2.3. Delete media from databases
			if err := s.mediaRepo.DeleteOne(ctx, mediaId); err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}
		}
	}

	// 3. Create Media and upload media
	if len(command.Media) > 0 {
		for _, file := range command.Media {
			// 3.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := media.SaveMedia(&file)

			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}

			// 3.2. create Media model and save to database
			mediaEntity, err := postEntity.NewMedia(postUpdated.ID, mediaUrl)
			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}

			_, err = s.mediaRepo.CreateOne(ctx, mediaEntity)
			if err != nil {
				return nil, response.NewServerFailedError(err.Error())
			}
		}
	}

	// 4. Delete cache post
	s.postCache.DeletePost(ctx, *command.PostId)

	return &postCommand.UpdatePostCommandResult{
		Post: mapper.NewPostResultFromEntity(postUpdated),
	}, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	command *postCommand.DeletePostCommand,
) (err error) {
	postFound, err := s.postRepo.GetById(ctx, *command.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return response.NewDataNotFoundError("post not found")
	}

	// 1. Get media array of post
	medias, err := s.mediaRepo.GetMany(ctx, "post_id=?", command.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 2. Delete media from database and folder
	for _, mediaRecord := range medias {
		// 2.1. Delete media from folder
		if err := media.DeleteMedia(mediaRecord.MediaUrl); err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// 2.1. Delete media from databases
		if err := s.mediaRepo.DeleteOne(ctx, mediaRecord.ID); err != nil {
			return response.NewServerFailedError(err.Error())
		}
	}

	deleteCondition := map[string]interface{}{
		"post_id": command.PostId,
	}

	// 3. Delete new feed
	err = s.newFeedRepo.DeleteMany(ctx, deleteCondition)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Delete advertise and bill
	err = s.advertiseRepo.DeleteMany(ctx, deleteCondition)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Delete post
	_, err = s.postRepo.DeleteOne(ctx, *command.PostId)
	if err != nil {
		return err
	}

	// 6. Delete post cache
	if err = s.deleteFeedCache(ctx, *command.PostId, postFound.UserId); err != nil {
		return err
	}
	return nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	query *postQuery.GetOnePostQuery,
) (result *postQuery.GetOnePostQueryResult, err error) {
	// 1. Get post
	postFound := s.postCache.GetPost(ctx, query.PostId)
	if postFound == nil {
		postFound, err = s.postRepo.GetById(ctx, query.PostId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
		if postFound == nil {
			return nil, response.NewDataNotFoundError("post not found")
		}
		go func(post *postEntity.Post) {
			s.postCache.SetPost(ctx, post)
		}(postFound)
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
				return nil, response.NewServerFailedError(err.Error())
			}
			if !isFriend {
				return nil, response.NewCustomError(response.ErrPostFriendAccess)
			}
		case consts.PRIVATE:
			return nil, response.NewCustomError(response.ErrPostPrivateAccess)
		default:
			return nil, response.NewCustomError(response.ErrPostPrivateAccess)
		}
	}

	// 2. Get user like
	postLikedQuery, err := postEntity.NewLikeUserPostEntity(query.AuthenticatedUserId, query.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}
	postLiked, err := s.likeUserPostRepo.CheckUserLikePost(ctx, postLikedQuery)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 3. Publish event to rabbitmq for statistic
	go func(postId uuid.UUID) {
		statisticEvent := common.StatisticEventResult{
			PostId:    postId,
			EventType: "clicks",
			Count:     1,
			Timestamp: time.Now(),
		}
		if err = s.postEventPublisher.PublishStatistic(ctx, statisticEvent, "stats.post"); err != nil {
			global.Logger.Error("Failed to publish statistic", zap.Error(err))
		}
	}(postFound.ID)

	// 4. Return
	return &postQuery.GetOnePostQueryResult{
		Post: mapper.NewPostWithLikedResultFromEntity(postFound, postLiked),
	}, nil
}

func (s *sPostUser) GetManyPosts(
	ctx context.Context,
	query *postQuery.GetManyPostQuery,
) (result *postQuery.GetManyPostQueryResult, err error) {
	// 1. Get post id list from cache
	postIDs, paging := s.postCache.GetFeeds(
		ctx, consts.RK_PERSONAL_POST, query.UserID, query.Limit, query.Page,
	)

	cacheFailed := false
	if len(postIDs) == 0 {
		cacheFailed = true
	}

	if query.Content != "" ||
		query.Location != "" ||
		query.IsAdvertisement != nil ||
		!query.CreatedAt.IsZero() ||
		query.SortBy != "" {
		cacheFailed = true
	}

	// 2. Cache hit
	var posts []*postEntity.Post
	if !cacheFailed {
		var wg sync.WaitGroup
		var postMap sync.Map
		cacheErrorOccurred := false

		for _, postID := range postIDs {
			wg.Add(1)
			go func(postID uuid.UUID) {
				defer wg.Done()
				post := s.postCache.GetPost(ctx, postID)
				if post == nil {
					post, err = s.postRepo.GetById(ctx, postID)
					if err != nil || post == nil {
						global.Logger.Warn("Failed to get post", zap.String("postId", postID.String()))
						cacheErrorOccurred = true
						if err = s.deleteFeedCache(ctx, postID, query.UserID); err != nil {
							return
						}
						return
					}
					s.postCache.SetPost(ctx, post)
				}
				postMap.Store(postID, post)
			}(postID)
		}
		wg.Wait()

		if cacheErrorOccurred {
			cacheFailed = true
		}

		if !cacheFailed {
			for _, postID := range postIDs {
				if post, ok := postMap.Load(postID); ok {
					posts = append(posts, post.(*postEntity.Post))
				}
			}
		}
	}

	// 3. Cache miss or cache handle error
	if cacheFailed {
		global.Logger.Warn("cache failed to get post, fallback to database")
		var pagingResp *response.PagingResponse
		posts, pagingResp, err = s.postRepo.GetMany(ctx, query)
		if err != nil {
			return nil, err
		}
		paging = pagingResp

		postIDs = make([]uuid.UUID, 0, len(posts))
		var wg sync.WaitGroup
		for _, post := range posts {
			postIDs = append(postIDs, post.ID)
			wg.Add(1)
			go func(p *postEntity.Post) {
				defer wg.Done()
				s.postCache.SetPost(ctx, p)
			}(post)
		}
		wg.Wait()

		s.postCache.SetFeeds(ctx, consts.RK_PERSONAL_POST, query.UserID, postIDs, pagingResp)
	}

	// 4. Get list user post like
	isLikedListQuery := &postQuery.CheckUserLikeManyPostQuery{
		PostIds:             postIDs,
		AuthenticatedUserId: query.AuthenticatedUserId,
	}
	isLikedList, err := s.likeUserPostRepo.CheckUserLikeManyPost(ctx, isLikedListQuery)
	if err != nil {
		return nil, err
	}

	// 5. Publish event to rabbitmq for statistic
	var wg sync.WaitGroup
	for _, post := range posts {
		postId := post.ID
		wg.Add(1)
		go func(postId uuid.UUID) {
			defer wg.Done()
			statisticEvent := common.StatisticEventResult{
				PostId:    postId,
				EventType: "impression",
				Count:     1,
				Timestamp: time.Now(),
			}
			if err = s.postEventPublisher.PublishStatistic(ctx, statisticEvent, "stats.post"); err != nil {
				global.Logger.Error("Failed to publish statistic", zap.Error(err))
			}
		}(postId)
	}
	wg.Wait()

	// Map to return
	var postResults []*common.PostResultWithLiked
	for _, post := range posts {
		postResult := mapper.NewPostWithLikedResultFromEntity(post, isLikedList[post.ID])
		postResults = append(postResults, postResult)
	}

	return &postQuery.GetManyPostQueryResult{
		Posts:          postResults,
		PagingResponse: paging,
	}, nil
}

func (s *sPostUser) GetTrendingPost(
	ctx context.Context,
	query *postQuery.GetTrendingPostQuery,
) (result *postQuery.GetManyPostQueryResult, err error) {
	// 1. Get Trending post
	posts, pagingResp, err := s.postRepo.GetTrendingPost(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Get list user post like
	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	isLikedListQuery := &postQuery.CheckUserLikeManyPostQuery{
		PostIds:             postIDs,
		AuthenticatedUserId: query.AuthenticatedUserId,
	}
	isLikedList, err := s.likeUserPostRepo.CheckUserLikeManyPost(ctx, isLikedListQuery)
	if err != nil {
		return nil, err
	}

	// Map to return
	var postResults []*common.PostResultWithLiked
	for _, post := range posts {
		postResult := mapper.NewPostWithLikedResultFromEntity(post, isLikedList[post.ID])
		postResults = append(postResults, postResult)
	}

	return &postQuery.GetManyPostQueryResult{
		Posts:          postResults,
		PagingResponse: pagingResp,
	}, nil
}

func (s *sPostUser) CheckPostOwner(
	ctx context.Context,
	query *postQuery.CheckPostOwnerQuery,
) (result *postQuery.CheckPostOwnerQueryResult, err error) {
	isOwner, err := s.postRepo.CheckPostOwner(ctx, query.PostId, query.UserId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &postQuery.CheckPostOwnerQueryResult{
		IsOwner: isOwner,
	}, nil
}

func (s *sPostUser) deleteFeedCache(ctx context.Context, postID, userID uuid.UUID) error {
	s.postCache.DeletePost(ctx, postID)
	s.postCache.DeleteFeeds(ctx, consts.RK_PERSONAL_POST, userID)
	s.postCache.DeleteFeeds(ctx, consts.RK_USER_FEED, userID)
	friends, err := s.friendRepo.GetFriendIds(ctx, userID)
	if err != nil {
		return err
	}
	if len(friends) == 0 {
		return nil
	}

	s.postCache.DeleteFriendFeeds(ctx, consts.RK_USER_FEED, friends)
	return nil
}
