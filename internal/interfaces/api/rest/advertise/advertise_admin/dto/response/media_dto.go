package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
)

type MediaDto struct {
	ID        uint      `json:"id"`
	PostId    uuid.UUID `json:"post_id"`
	MediaUrl  string    `json:"media_url"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToMediaDto(mediaResult []*common.MediaResult) []*MediaDto {
	var mediaDtos []*MediaDto
	for _, media := range mediaResult {
		mediaDto := &MediaDto{
			ID:        media.ID,
			PostId:    media.PostId,
			MediaUrl:  media.MediaUrl,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		}
		mediaDtos = append(mediaDtos, mediaDto)
	}
	return mediaDtos
}
