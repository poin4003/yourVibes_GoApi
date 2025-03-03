package implement

import (
	"context"

	mediaQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/media"
)

type sMedia struct{}

func NewMediaImplement() *sMedia {
	return &sMedia{}
}

func (s *sMedia) GetMedia(
	ctx context.Context,
	query *mediaQuery.MediaQuery,
) (result *mediaQuery.MediaQueryResult, err error) {
	// 1. Get file path
	filePath, err := media.GetMedia(query.FileName)
	if err != nil {
		return nil, err
	}

	return &mediaQuery.MediaQueryResult{
		FilePath: filePath,
	}, nil
}
