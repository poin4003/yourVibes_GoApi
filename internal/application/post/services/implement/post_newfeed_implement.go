package implement

import (
	"context"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type sPostNewFeed struct {
	userRepo         post_repo.IUserRepository
	postRepo         post_repo.IPostRepository
	likeUserPostRepo post_repo.ILikeUserPostRepository
	newFeedRepo      post_repo.INewFeedRepository
}

func NewPostNewFeedImplement(
	userRepo post_repo.IUserRepository,
	postRepo post_repo.IPostRepository,
	likeUserPostRepo post_repo.ILikeUserPostRepository,
	newFeedRepo post_repo.INewFeedRepository,
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
	command *post_command.DeleteNewFeedCommand,
) (result *post_command.DeleteNewFeedCommandResult, err error) {
	result = &post_command.DeleteNewFeedCommandResult{}

	err = s.newFeedRepo.DeleteOne(ctx, command.UserId, command.PostId)
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sPostNewFeed) GetNewFeeds(
	ctx context.Context,
	query *post_query.GetNewFeedQuery,
) (result *post_query.GetNewFeedResult, err error) {
	result = &post_query.GetNewFeedResult{}

	postEntities, paging, err := s.newFeedRepo.GetMany(ctx, query)
	if err != nil {
		result.Posts = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
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
