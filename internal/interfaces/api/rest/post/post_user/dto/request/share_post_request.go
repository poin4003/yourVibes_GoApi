package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type SharePostRequest struct {
	Content  string              `form:"content,omitempty" binding:"omitempty"`
	Privacy  consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty"`
	Location string              `form:"location,omitempty"`
}

func ValidateSharePostRequest(req interface{}) error {
	dto, ok := req.(*SharePostRequest)
	if !ok {
		return fmt.Errorf("validate SharePostRequest failed")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Content, validation.Length(2, 1000)),
		validation.Field(&dto.Privacy, validation.In(consts.PrivacyLevels...)),
	)
}

func (req *SharePostRequest) ToSharePostCommand(
	postId uuid.UUID,
	userId uuid.UUID,
) (*postCommand.SharePostCommand, error) {
	return &postCommand.SharePostCommand{
		PostId:   postId,
		UserId:   userId,
		Content:  req.Content,
		Privacy:  req.Privacy,
		Location: req.Location,
	}, nil
}
