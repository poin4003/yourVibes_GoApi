package implement

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"go.uber.org/zap"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sPostNewFeed struct {
	userRepo           postRepo.IUserRepository
	postRepo           postRepo.IPostRepository
	likeUserPostRepo   postRepo.ILikeUserPostRepository
	newFeedRepo        postRepo.INewFeedRepository
	postCache          cache.IPostCache
	postEventPublisher *producer.PostEventPublisher
}

func NewPostNewFeedImplement(
	userRepo postRepo.IUserRepository,
	postRepo postRepo.IPostRepository,
	likeUserPostRepo postRepo.ILikeUserPostRepository,
	newFeedRepo postRepo.INewFeedRepository,
	postCache cache.IPostCache,
	postEventPublisher *producer.PostEventPublisher,
) *sPostNewFeed {
	return &sPostNewFeed{
		userRepo:           userRepo,
		postRepo:           postRepo,
		likeUserPostRepo:   likeUserPostRepo,
		newFeedRepo:        newFeedRepo,
		postCache:          postCache,
		postEventPublisher: postEventPublisher,
	}
}

func (s *sPostNewFeed) DeleteNewFeed(
	ctx context.Context,
	command *postCommand.DeleteNewFeedCommand,
) (err error) {
	err = s.newFeedRepo.DeleteOne(ctx, command.UserId, command.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	s.postCache.DeleteFeeds(ctx, consts.RK_USER_FEED, command.UserId)

	return nil
}

func (s *sPostNewFeed) GetNewFeeds(
	ctx context.Context,
	query *postQuery.GetNewFeedQuery,
) (result *postQuery.GetNewFeedResult, err error) {
	// 1. Get post id list from cache
	postIDs, paging := s.postCache.GetFeeds(
		ctx, consts.RK_USER_FEED, query.UserId, query.Limit, query.Page,
	)

	// 2. Cache miss
	var posts []*postEntity.Post
	if len(postIDs) == 0 {
		var postEntities []*postEntity.Post
		var pagingResp *response.PagingResponse
		postEntities, pagingResp, err = s.newFeedRepo.GetMany(ctx, query)
		if err != nil {
			return nil, err
		}
		posts = postEntities
		paging = pagingResp
		postIDs = make([]uuid.UUID, 0, len(postEntities))
		var wg sync.WaitGroup
		for _, post := range postEntities {
			postIDs = append(postIDs, post.ID)
			wg.Add(1)
			go func(p *postEntity.Post) {
				defer wg.Done()
				s.postCache.SetPost(ctx, p)
			}(post)
		}
		wg.Wait()

		s.postCache.SetFeeds(ctx, consts.RK_USER_FEED, query.UserId, postIDs, pagingResp)
	} else {
		// Cache hit
		posts = make([]*postEntity.Post, 0, len(postIDs))
		var wg sync.WaitGroup
		ch := make(chan *postEntity.Post, len(postIDs))

		for _, postID := range postIDs {
			wg.Add(1)
			go func(postID uuid.UUID) {
				defer wg.Done()
				post := s.postCache.GetPost(ctx, postID)
				if post == nil {
					post, err = s.postRepo.GetById(ctx, postID)
					if err == nil || post != nil {
						s.postCache.SetPost(ctx, post)
					}
				}
				ch <- post
			}(postID)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()

		for post := range ch {
			posts = append(posts, post)
		}
	}

	// 3. Get list user post like
	isLikedListQuery := &postQuery.CheckUserLikeManyPostQuery{
		PostIds:             postIDs,
		AuthenticatedUserId: query.UserId,
	}
	isLikedList, err := s.likeUserPostRepo.CheckUserLikeManyPost(ctx, isLikedListQuery)
	if err != nil {
		return nil, err
	}

	// Publish event to rabbitmq for statistic
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

	var postResults []*common.PostResultWithLiked
	for _, post := range posts {
		postResult := mapper.NewPostWithLikedResultFromEntity(post, isLikedList[post.ID])
		postResults = append(postResults, postResult)
	}

	return &postQuery.GetNewFeedResult{
		Posts:          postResults,
		PagingResponse: paging,
	}, nil
}
