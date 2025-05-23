package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
	reportEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
)

func NewMediaResult(
	media []*reportEntity.Media,
) []*common.MediaResult {
	if media == nil {
		return nil
	}

	var mediaResults []*common.MediaResult
	for _, mediaEntity := range media {
		mediaResult := &common.MediaResult{
			ID:        mediaEntity.ID,
			PostId:    mediaEntity.PostId,
			MediaUrl:  mediaEntity.MediaUrl,
			Status:    mediaEntity.Status,
			CreatedAt: mediaEntity.CreatedAt,
			UpdatedAt: mediaEntity.UpdatedAt,
		}
		mediaResults = append(mediaResults, mediaResult)
	}

	return mediaResults
}
