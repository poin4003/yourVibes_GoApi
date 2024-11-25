package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	post_command "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
	"strings"
)

type UpdatePostRequest struct {
	Content  *string                `form:"content,omitempty"`
	Privacy  *consts.PrivacyLevel   `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Location *string                `form:"location,omitempty"`
	MediaIDs []uint                 `form:"media_ids,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"omitempty,files"`
}

func ValidateUpdatePostRequest(req interface{}) error {
	dto, ok := req.(*UpdatePostRequest)
	if !ok {
		return fmt.Errorf("validate CreatePostRequest failed")
	}

	if dto.Media != nil && len(dto.Media) > 0 {
		for _, fileHeader := range dto.Media {
			if err := validateMediaForUpdate(&fileHeader); err != nil {
				return err
			}
		}
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Content, validation.Min(2)),
		validation.Field(&dto.Privacy, validation.In(consts.PRIVATE, consts.PUBLIC, consts.FRIEND_ONLY)),
	)
}

func validateMediaForUpdate(value interface{}) error {
	fileHeader, ok := value.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("invalid file format")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") || !strings.HasPrefix(contentType, "video/") {
		return fmt.Errorf("file must be an image or video")
	}

	if fileHeader.Size > 10*1024*1024 {
		return fmt.Errorf("file size must be less than 10M")
	}

	return nil
}

func (req *UpdatePostRequest) ToUpdatePostCommand(
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
