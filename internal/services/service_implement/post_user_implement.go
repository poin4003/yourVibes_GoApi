package service_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
)

type sPostUser struct {
	userRepo  repository.IUserRepository
	postRepo  repository.IPostRepository
	mediaRepo repository.IMediaRepository
}

func NewPostUserImplement(
	userRepo repository.IUserRepository,
	postRepo repository.IPostRepository,
	mediaRepo repository.IMediaRepository,
) *sPostUser {
	return &sPostUser{
		userRepo:  userRepo,
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	postModel *model.Post,
	inMedia []multipart.File,
) (post *model.Post, resultCode int, err error) {
	// 1. CreatePost
	newPost, err := s.postRepo.CreatePost(ctx, postModel)
	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	// 2. Create Media and upload media to cloudinary_util
	if len(inMedia) > 0 {
		for _, file := range inMedia {
			// 2.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			// 2.2. create Media model and save to database
			mediaTemp := &model.Media{
				PostId:   newPost.ID,
				MediaUrl: mediaUrl,
			}

			_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
			if err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	userFound, err := s.userRepo.GetUser(ctx, "id=?", postModel.UserId)
	if err != nil {
		return nil, response.ErrDataNotFound, fmt.Errorf("failed to get user: %w", err)
	}

	userFound.PostCount++

	_, err = s.userRepo.UpdateUser(ctx, userFound.ID, map[string]interface{}{
		"post_count": userFound.PostCount,
	})

	return newPost, response.ErrCodeSuccess, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	postId uuid.UUID,
	updateData map[string]interface{},
	deleteMediaIds []uint,
	inMedia []multipart.File,
) (post *model.Post, resultCode int, err error) {
	// 1. update post information
	postModel, err := s.postRepo.UpdatePost(ctx, postId, updateData)
	if err != nil {
		return nil, response.ErrServerFailed, err
	}

	// 2. delete media in database and delete media from cloudinary
	if len(deleteMediaIds) > 0 {
		for _, mediaId := range deleteMediaIds {
			// 2.1. Get media information from database
			media, err := s.mediaRepo.GetMedia(ctx, "id=?", mediaId)
			if err != nil {
				return nil, response.ErrDataNotFound, fmt.Errorf("failed to get media record: %w", err)
			}

			// 2.2. Delete media from cloudinary
			if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to delete media record: %w", err)
			}

			// 2.3. Delete media from databases
			if err := s.mediaRepo.DeleteMedia(ctx, mediaId); err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to delete media record: %w", err)
			}
		}
	}

	fmt.Println(len(inMedia))
	// 3. Create Media and upload media to cloudinary_util
	if len(inMedia) > 0 {
		for _, file := range inMedia {
			// 3.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			// 3.2. create Media model and save to database
			mediaTemp := &model.Media{
				PostId:   postId,
				MediaUrl: mediaUrl,
			}

			_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
			if err != nil {
				return nil, response.ErrServerFailed, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}

	return postModel, response.ErrCodeSuccess, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	postId uuid.UUID,
) (resultCode int, err error) {
	// 1. Get media array of post
	medias, err := s.mediaRepo.GetManyMedia(ctx, "post_id=?", postId)
	if err != nil {
		return response.ErrDataNotFound, fmt.Errorf("failed to get media records: %w", err)
	}

	// 2. Delete media from database and cloudinary
	for _, media := range medias {
		// 2.1. Delete media from cloudinary
		if err := cloudinary_util.DeleteMediaFromCloudinary(media.MediaUrl); err != nil {
			return response.ErrServerFailed, fmt.Errorf("failed to delete media record: %w", err)
		}

		// 2.1. Delete media from databases
		if err := s.mediaRepo.DeleteMedia(ctx, media.ID); err != nil {
			return response.ErrServerFailed, fmt.Errorf("failed to delete media record: %w", err)
		}
	}

	deletePostErr := s.postRepo.DeletePost(ctx, postId)

	if deletePostErr != nil {
		return response.ErrServerFailed, fmt.Errorf(deletePostErr.Error())
	}

	return response.ErrCodeSuccess, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	postId uuid.UUID,
) (post *model.Post, resultCode int, err error) {
	postModel, err := s.postRepo.GetPost(ctx, "id=?", postId)
	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	return postModel, response.ErrCodeSuccess, nil
}

func (s *sPostUser) GetManyPosts(
	ctx context.Context,
	query *query_object.PostQueryObject,
) (posts []*model.Post, resultCode int, err error) {
	postModels, err := s.postRepo.GetManyPost(ctx, query)

	if err != nil {
		return nil, response.ErrDataNotFound, err
	}

	return postModels, response.ErrCodeSuccess, nil
}
