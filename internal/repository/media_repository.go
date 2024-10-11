package repository

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

type (
	IMediaRepository interface {
		CreateMedia(ctx context.Context, media *model.Media) (*model.Media, error)
		UpdateMedia(ctx context.Context, mediaId uint, updateData map[string]interface{}) (*model.Media, error)
		DeleteMedia(ctx context.Context, mediaId uint) error
		GetMedia(ctx context.Context, query interface{}, args ...interface{}) (*model.Media, error)
		GetManyMedia(ctx context.Context) ([]*model.Media, error)
	}
)

var (
	localMedia IMediaRepository
)

func Media() IMediaRepository {
	if localMedia == nil {
		panic("repository_implement localMedia not found for interface IMedia")
	}

	return localMedia
}

func InitMediaRepository(i IMediaRepository) {
	localMedia = i
}
