package service_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type sPostShare struct {
	userRepo  repository.IUserRepository
	postRepo  repository.IPostRepository
	mediaRepo repository.IMediaRepository
}

func NewPostShareImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	mediaRepo repository.IMediaRepository,
) *sPostShare {
	return &sPostShare{
		userRepo:  userRepo,
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostShare) SharePost(
	ctx context.Context,
	postId uuid.UUID,
	userId uuid.UUID,
) (post *model.Post, resultCode int, err error) {
	// 1. Find root post by post_id
	postModel, err := s.postRepo.GetPost(ctx, "id = ?", postId)
	if err != nil {
		return nil, response.ErrServerFailed, err
	}

	if postModel == nil {
		return nil, response.ErrDataNotFound, nil
	}

	// 2. Create new post (parent_id = post_id, user_id = userId)
	if postModel.ParentId == nil {
		// 2.1. Copy post info from root post
		newPost := &model.Post{
			UserId:   userId,
			ParentId: &postModel.ID,
			Content:  postModel.Content,
			Location: postModel.Location,
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreatePost(ctx, newPost)
		if err != nil {
			return nil, response.ErrServerFailed, err
		}

		// 2.3. Copy media to new share post from root post
		if len(postModel.Media) > 0 {
			for _, media := range postModel.Media {
				// 2.3.1. Copy media from root post
				mediaTemp := &model.Media{
					PostId:   newSharePost.ID,
					MediaUrl: media.MediaUrl,
				}

				// 2.3.2. Create media for new post
				_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
				if err != nil {
					return nil, response.ErrServerFailed, fmt.Errorf("failed to create media record: %w", err)
				}
			}
		}
		return newSharePost, response.ErrCodeSuccess, nil
	} else {
		// 3. Find actually root post
		rootPost, err := s.postRepo.GetPost(ctx, "id=?", postModel.ParentId)
		if err != nil {
			return nil, response.ErrServerFailed, err
		}

		if rootPost == nil {
			return nil, response.ErrDataNotFound, nil
		}

		// 3.1. Copy post info from root post
		newPost := &model.Post{
			UserId:   userId,
			ParentId: &rootPost.ID,
			Content:  rootPost.Content,
			Location: rootPost.Location,
		}

		// 2.2. Create new post
		newSharePost, err := s.postRepo.CreatePost(ctx, newPost)
		if err != nil {
			return nil, response.ErrServerFailed, err
		}

		// 2.3. Copy media to new share post from root post
		if len(rootPost.Media) > 0 {
			for _, media := range rootPost.Media {
				// 2.3.1. Copy media from root post
				mediaTemp := &model.Media{
					PostId:   newSharePost.ID,
					MediaUrl: media.MediaUrl,
				}

				// 2.3.2. Create media for new post
				_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
				if err != nil {
					return nil, response.ErrServerFailed, fmt.Errorf("failed to create media record: %w", err)
				}
			}
		}
		return newSharePost, response.ErrCodeSuccess, nil
	}
}
