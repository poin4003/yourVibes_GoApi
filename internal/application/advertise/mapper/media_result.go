package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	advertiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
)

func NewMediaResult(
	media []*advertiseEntity.Media,
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
