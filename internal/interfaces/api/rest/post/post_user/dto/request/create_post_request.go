package request

import (
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

type CreatePostInput struct {
	Content  string                 `form:"content,omitempty" binding:"omitempty"`
	Privacy  consts.PrivacyLevel    `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Location string                 `form:"location,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"omitempty,files"`
}

func (req *CreatePostInput) ToCreatePostCommand(
	userId uuid.UUID,
	media []multipart.File,
) (*post_command.CreatePostCommand, error) {
	return &post_command.CreatePostCommand{
		UserId:   userId,
		Content:  req.Content,
		Privacy:  req.Privacy,
		Location: req.Location,
		Media:    media,
	}, nil
}
