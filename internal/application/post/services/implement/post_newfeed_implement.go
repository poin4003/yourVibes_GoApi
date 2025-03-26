package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"go.uber.org/zap"
	"sync"
	"time"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sPostNewFeed struct {
	userRepo           postRepo.IUserRepository
	postRepo           postRepo.IPostRepository
	likeUserPostRepo   postRepo.ILikeUserPostRepository
	newFeedRepo        postRepo.INewFeedRepository
	postEventPublisher *producer.PostEventPublisher
}

func NewPostNewFeedImplement(
	userRepo postRepo.IUserRepository,
	postRepo postRepo.IPostRepository,
	likeUserPostRepo postRepo.ILikeUserPostRepository,
	newFeedRepo postRepo.INewFeedRepository,
	postEventPublisher *producer.PostEventPublisher,
) *sPostNewFeed {
	return &sPostNewFeed{
		userRepo:           userRepo,
		postRepo:           postRepo,
		likeUserPostRepo:   likeUserPostRepo,
		newFeedRepo:        newFeedRepo,
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

	return nil
}

func (s *sPostNewFeed) GetNewFeeds(
	ctx context.Context,
	query *postQuery.GetNewFeedQuery,
) (result *postQuery.GetNewFeedResult, err error) {
	postEntities, paging, err := s.newFeedRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var wg sync.WaitGroup
	for _, post := range postEntities {
		wg.Add(1)
		go func(post *postEntity.PostWithLiked) {
			defer wg.Done()
			statisticEvent := common.StatisticEventResult{
				PostId:    post.ID,
				EventType: "impression",
				Count:     1,
				Timestamp: time.Now(),
			}
			if err = s.postEventPublisher.PublishStatistic(ctx, statisticEvent, "stats.post"); err != nil {
				global.Logger.Error("Failed to publish statistic", zap.Error(err))
			}
		}(post)
	}
	wg.Wait()

	var postResults []*common.PostResultWithLiked
	for _, postEntity := range postEntities {
		postResults = append(postResults, mapper.NewPostWithLikedResultFromEntity(postEntity))
	}

	return &postQuery.GetNewFeedResult{
		Posts:          postResults,
		PagingResponse: paging,
	}, nil
}
