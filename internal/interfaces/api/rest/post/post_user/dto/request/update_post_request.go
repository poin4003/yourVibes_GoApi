package request

import (
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

type UpdatePostInput struct {
	Content  *string                `form:"content,omitempty"`
	Privacy  *consts.PrivacyLevel   `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Location *string                `form:"location,omitempty"`
	MediaIDs []uint                 `form:"media_ids,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"omitempty,files"`
}

func (req *UpdatePostInput) ToUpdatePostCommand(
	postId *uuid.UUID,
	media []multipart.File,
) (*post_command.UpdatePostCommand, error) {
	return &post_command.UpdatePostCommand{
		PostId:   postId,
		Content:  req.Content,
		Privacy:  req.Privacy,
		Location: req.Location,
		MediaIDs: req.MediaIDs,
		Media:    media,
	}, nil
}
