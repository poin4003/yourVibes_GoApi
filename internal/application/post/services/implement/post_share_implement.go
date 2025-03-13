package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	postValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/validator"
	postRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
			return nil, response.NewServerFailedError(err.Error())
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		validatePost, err := postValidator.NewValidatedPost(newSharePost)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		result.Post = mapper.NewPostResultFromValidateEntity(validatePost)

		return result, nil
	} else {
		// 3. Find actually root post
		rootPost, err := s.postRepo.GetById(ctx, *postFound.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if rootPost == nil {
			return nil, response.NewDataNotFoundError("root post not found")
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
			return nil, response.NewServerFailedError(err.Error())
		}
		// 3.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		validatePost, err := postValidator.NewValidatedPost(newSharePost)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		result.Post = mapper.NewPostResultFromValidateEntity(validatePost)

		return result, nil
	}
}
