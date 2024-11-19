package implement

import (
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sUserNewFeed struct {
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
) *sUserNewFeed {
	return &sUserNewFeed{
		userRepo:         userRepo,
		postRepo:         postRepo,
		likeUserPostRepo: likeUserPostRepo,
		newFeedRepo:      newFeedRepo,
	}
}

//func (s *sUserNewFeed) DeleteNewFeed(
//	ctx context.Context,
//	userId uuid.UUID,
//	postId uuid.UUID,
//) (resultCode int, httpStatusCode int, err error) {
//	err = s.newFeedRepo.DeleteNewFeed(ctx, userId, postId)
//	if err != nil {
//		return response.ErrServerFailed, http.StatusInternalServerError, err
//	}
//
//	return response.ErrCodeSuccess, http.StatusOK, nil
//}
//
//func (s *sUserNewFeed) GetNewFeeds(
//	ctx context.Context,
//	userId uuid.UUID,
//	query *query.NewFeedQueryObject,
//) (postDtos []*response2.PostDto, pagingResponse *response.PagingResponse, resultCode int, httpStatusCode int, err error) {
//	postModels, paging, err := s.newFeedRepo.GetManyNewFeed(ctx, userId, query)
//	if err != nil {
//		return nil, nil, response.ErrServerFailed, http.StatusInternalServerError, err
//	}
//
//	for _, post := range postModels {
//		isLiked, _ := s.likeUserPostRepo.CheckUserLikePost(ctx, &models.LikeUserPost{
//			PostId: post.ID,
//			UserId: userId,
//		})
//
//		postDto := mapper.MapPostToPostDto(post, isLiked)
//		postDtos = append(postDtos, postDto)
//	}
//
//	return postDtos, paging, response.ErrCodeSuccess, http.StatusOK, nil
//}
