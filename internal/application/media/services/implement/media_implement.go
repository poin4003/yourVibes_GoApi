package implement

import (
	"context"
	media_query "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/media"
	"net/http"
)

type sMedia struct{}

func NewMediaImplement() *sMedia {
	return &sMedia{}
}

func (s *sMedia) GetMedia(
	ctx context.Context,
	query *media_query.MediaQuery,
) (result *media_query.MediaQueryResult, err error) {
	result = &media_query.MediaQueryResult{}
	result.FilePath = ""
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
