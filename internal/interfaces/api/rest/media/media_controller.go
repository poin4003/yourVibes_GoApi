package media

import (
	"github.com/gin-gonic/gin"
	mediaServiceQuery "github.com/poin4003/yourVibes_GoApi/internal/application/media/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/media/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cMedia struct{}

func NewMediaController() *cMedia {
	return &cMedia{}
}

func (c *cMedia) GetMedia(ctx *gin.Context) {
	// 1. Get file name from path
	fileName := ctx.Param("file_name")

	// 2. Call service to get file path
	query := &mediaServiceQuery.MediaQuery{
		FileName: fileName,
	}
	result, err := services.Media().GetMedia(ctx, query)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	ctx.File(result.FilePath)
}
