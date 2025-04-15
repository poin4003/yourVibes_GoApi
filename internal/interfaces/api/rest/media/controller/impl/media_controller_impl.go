package impl

import (
	"github.com/gin-gonic/gin"
	mediaQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/media/services"
	"net/http"
)

type cMedia struct {
	mediaService services.IMedia
}

func NewMediaController(
	mediaService services.IMedia,
) *cMedia {
	return &cMedia{
		mediaService: mediaService,
	}
}

func (c *cMedia) GetMedia(ctx *gin.Context) {
	// 1. Get file name from path
	fileName := ctx.Param("file_name")
	rangeHeader := ctx.Request.Header.Get("Range")

	// 2. Call service to handle media streaming
	query := &mediaQuery.MediaQuery{
		FileName:    fileName,
		RangeHeader: rangeHeader,
	}
	result, err := c.mediaService.GetMedia(ctx, query)
	if err != nil {
		ctx.Error(err)
		return
	}

	defer result.File.Close()

	// 3. Set headers and serve content
	for key, value := range result.Headers {
		ctx.Header(key, value)
	}
	if result.StatusCode != 0 {
		ctx.Status(result.StatusCode)
	}
	http.ServeContent(ctx.Writer, ctx.Request, fileName, result.ModTime, result.File)
}
