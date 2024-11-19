package request

import (
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type SharePostInput struct {
	Content  string              `form:"content,omitempty" binding:"omitempty"`
	Privacy  consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Location string              `form:"location,omitempty"`
}

func (req *SharePostInput) ToSharePostCommand(
	postId uuid.UUID,
	userId uuid.UUID,
) (*post_command.SharePostCommand, error) {
	return &post_command.SharePostCommand{
		PostId:   postId,
		UserId:   userId,
		Content:  req.Content,
		Privacy:  req.Privacy,
		Location: req.Location,
	}, nil
}
