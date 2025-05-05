package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/media"
	"io"
	"net/http"
	"os"
	"strconv"

	mediaQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
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
	fileSize := fileInfo.Size()

	// 4. Prepare result
	result = &mediaQuery.MediaQueryResult{
		RawFile: file,
		ModTime: fileInfo.ModTime(),
		Headers: make(map[string]string),
	}

	// 5. Handle range header for streaming
	if query.RangeHeader == "" {
		query.RangeHeader = "bytes=0-" + strconv.FormatInt(fileSize-1, 10)
	}

	start, end, err := media.ParseRange(query.RangeHeader, fileSize)
	if err != nil {
		file.Close()
		return nil, response.NewServerFailedError(err.Error())
	}

	result.Headers["Accept-Ranges"] = "bytes"
	result.Headers["Content-Type"] = "video/mp4"
	result.Headers["Content-Length"] = strconv.FormatInt(end-start+1, 10)
	result.Headers["Content-Range"] = "bytes " + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10) + "/" + strconv.FormatInt(fileSize, 10)
	result.StatusCode = http.StatusPartialContent

	result.File = io.NewSectionReader(file, start, end-start+1)

	return result, nil
}
