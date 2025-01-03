package implement

import (
	"context"
	"net/http"

	mediaQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
	result = &mediaQuery.MediaQueryResult{
		FilePath:       "",
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get file path
	filePath, err := media.GetMedia(query.FileName)
	if err != nil {
		return nil, err
	}

	result.FilePath = filePath
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
