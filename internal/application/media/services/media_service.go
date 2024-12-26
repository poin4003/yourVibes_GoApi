package services

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
)

type (
	IMedia interface {
		GetMedia(ctx context.Context, query *query.MediaQuery) (result *query.MediaQueryResult, err error)
	}
)

var (
	localMedia IMedia
)

func Media() IMedia {
	if localMedia == nil {
		panic("media implement localMedia not found for interface IMedia")
	}

	return localMedia
}

func InitMedia(i IMedia) {
	localMedia = i
}
