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
