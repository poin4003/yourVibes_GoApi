package implement

import (
	"context"
	"errors"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sPostShare struct {
	userRepo  postRepo.IUserRepository
	postRepo  postRepo.IPostRepository
	mediaRepo postRepo.IMediaRepository
}

func NewPostShareImplement(
	userRepo postRepo.IUserRepository,
	postRepo postRepo.IPostRepository,
	mediaRepo postRepo.IMediaRepository,
) *sPostShare {
	return &sPostShare{
		userRepo:  userRepo,
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostShare) SharePost(
	ctx context.Context,
	command *postCommand.SharePostCommand,
) (result *postCommand.SharePostCommandResult, err error) {
	result = &postCommand.SharePostCommandResult{}
	result.Post = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	// 1. Find post by post_id
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Create new post (parent_id = post_id, user_id = userId)
	if postFound.ParentId == nil {
		// 2.1. Copy post info from root post
		newPost, err := postEntity.NewPostForShare(
			command.UserId,
			command.Content,
			command.Privacy,
			command.Location,
			&command.PostId,
		)
		if err != nil {
			return result, err
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			return result, err
		}

		validatePost, err := postValidator.NewValidatedPost(newSharePost)
		if err != nil {
			return result, err
		}

		result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	} else {
		// 3. Find actually root post
		rootPost, err := s.postRepo.GetById(ctx, *postFound.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Post = nil
				result.ResultCode = response.ErrDataNotFound
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
			return result, err
		}

		// 3.1. Copy post info from root post
		newPost, err := postEntity.NewPostForShare(
			command.UserId,
			command.Content,
			command.Privacy,
			command.Location,
			&rootPost.ID,
		)
		if err != nil {
			return result, err
		}
		// 3.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			return result, err
		}

		validatePost, err := postValidator.NewValidatedPost(newSharePost)
		if err != nil {
			return result, err
		}

		result.Post = mapper.NewPostResultFromValidateEntity(validatePost)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}
}
