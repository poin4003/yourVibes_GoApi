package implement

import (
	"context"
	"net/http"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type sPostNewFeed struct {
	userRepo         postRepo.IUserRepository
	postRepo         postRepo.IPostRepository
	likeUserPostRepo postRepo.ILikeUserPostRepository
	newFeedRepo      postRepo.INewFeedRepository
}

func NewPostNewFeedImplement(
	userRepo postRepo.IUserRepository,
	postRepo postRepo.IPostRepository,
	likeUserPostRepo postRepo.ILikeUserPostRepository,
	newFeedRepo postRepo.INewFeedRepository,
) *sPostNewFeed {
	return &sPostNewFeed{
		userRepo:         userRepo,
		postRepo:         postRepo,
		likeUserPostRepo: likeUserPostRepo,
		newFeedRepo:      newFeedRepo,
	}
}

func (s *sPostNewFeed) DeleteNewFeed(
	ctx context.Context,
	command *postCommand.DeleteNewFeedCommand,
) (result *postCommand.DeleteNewFeedCommandResult, err error) {
	result = &postCommand.DeleteNewFeedCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}

	err = s.newFeedRepo.DeleteOne(ctx, command.UserId, command.PostId)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostNewFeed) GetNewFeeds(
	ctx context.Context,
	query *postQuery.GetNewFeedQuery,
) (result *postQuery.GetNewFeedResult, err error) {
	result = &postQuery.GetNewFeedResult{
		Posts:          nil,
		PagingResponse: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}

	postEntities, paging, err := s.newFeedRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var postResults []*common.PostResultWithLiked
	for _, postEntity := range postEntities {
		postResults = append(postResults, mapper.NewPostWithLikedResultFromEntity(postEntity))
	}

	result.Posts = postResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}
