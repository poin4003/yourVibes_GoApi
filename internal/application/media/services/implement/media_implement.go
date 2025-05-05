package implement

import (
	"context"
	mediaQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/media"
	"os"
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

	// 2. Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, response.NewDataNotFoundError(err.Error())
	}

	// 3. Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, response.NewServerFailedError(err.Error())
	}

	return &mediaQuery.MediaQueryResult{
		File:    file,
		ModTime: fileInfo.ModTime(),
	}, nil
}
