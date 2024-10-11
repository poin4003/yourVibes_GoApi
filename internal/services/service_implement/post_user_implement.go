package service_implement

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/cloudinary_util"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
	"mime/multipart"
	"net/http"
)

type sPostUser struct {
	postRepo  repository.IPostRepository
	mediaRepo repository.IMediaRepository
}

func NewPostUserImplement(
	postRepo repository.IPostRepository,
	mediaRepo repository.IMediaRepository,
) *sPostUser {
	return &sPostUser{
		postRepo:  postRepo,
		mediaRepo: mediaRepo,
	}
}

func (s *sPostUser) CreatePost(
	ctx context.Context,
	inPostData *vo.CreatePostInput,
	inMedia []multipart.File,
	userId uuid.UUID,
) (post *model.Post, resultCode int, err error) {
	// 1. CreatePost
	postTemp := &model.Post{
		UserId:   userId,
		Title:    inPostData.Title,
		Content:  inPostData.Content,
		Privacy:  inPostData.Privacy,
		Location: inPostData.Location,
	}

	newPost, err := s.postRepo.CreatePost(ctx, postTemp)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 2. Create Media and upload media to cloudinary_util
	if len(inMedia) > 0 {
		for _, file := range inMedia {
			// 2.1. upload to cloudinary and get mediaUrl
			mediaUrl, err := cloudinary_util.UploadMediaToCloudinary(file)

			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("failed to upload media to cloudinary: %w", err)
			}

			// 2.2. create Media model and save to database
			mediaTemp := &model.Media{
				PostId:   newPost.ID,
				MediaUrl: mediaUrl,
			}

			_, err = s.mediaRepo.CreateMedia(ctx, mediaTemp)
			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("failed to create media record: %w", err)
			}
		}
	}
	return newPost, http.StatusCreated, nil
}

func (s *sPostUser) UpdatePost(
	ctx context.Context,
	in *vo.UpdatePostInput,
) (post *model.Post, resultCode int, err error) {
	return &model.Post{}, 0, nil
}

func (s *sPostUser) DeletePost(
	ctx context.Context,
	email string) (resultCode int, err error) {
	return 0, nil
}

func (s *sPostUser) GetPost(
	ctx context.Context,
	query interface{}, args ...interface{},
) (post *model.Post, resultCode int, err error) {
	return &model.Post{}, 0, nil
}

func (s *sPostUser) GetAllPost(
	ctx context.Context,
) (posts []*model.Post, resultCode int, err error) {
	return []*model.Post{}, 0, nil
}
