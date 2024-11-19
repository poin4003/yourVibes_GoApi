package implement

import (
	"context"
	"errors"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/mapper"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	post_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"net/http"
)

type sPostShare struct {
	userRepo  post_repo.IUserRepository
	postRepo  post_repo.IPostRepository
	mediaRepo post_repo.IMediaRepository
}

func NewPostShareImplement(
	userRepo post_repo.IUserRepository,
	postRepo post_repo.IPostRepository,
	mediaRepo post_repo.IMediaRepository,
) *sPostShare {
	return &sPostShare{
		userRepo:  userRepo,
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostShare) SharePost(
	ctx context.Context,
	command *post_command.SharePostCommand,
) (result *post_command.SharePostCommandResult, err error) {
	result = &post_command.SharePostCommandResult{}
	// 1. Find post by post_id
	postEntity, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Post = nil
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.Post = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Create new post (parent_id = post_id, user_id = userId)
	if postEntity.ParentId == nil {
		// 2.1. Copy post info from root post
		newPost, err := post_entity.NewPostForShare(
			command.UserId,
			command.Content,
			command.Privacy,
			command.Location,
			&command.PostId,
		)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}

		result.Post = mapper.NewPostResultFromEntity(newSharePost)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	} else {
		// 3. Find actually root post
		rootPost, err := s.postRepo.GetById(ctx, *postEntity.ParentId)
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
		newPost, err := post_entity.NewPostForShare(
			command.UserId,
			command.Content,
			command.Privacy,
			command.Location,
			&rootPost.ID,
		)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}
		// 3.2. Create new post
		newSharePost, err := s.postRepo.CreateOne(ctx, newPost)
		if err != nil {
			result.Post = nil
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
			return result, err
		}

		result.Post = mapper.NewPostResultFromEntity(newSharePost)
		result.ResultCode = response.ErrCodeSuccess
		result.HttpStatusCode = http.StatusOK
		return result, nil
	}
}
