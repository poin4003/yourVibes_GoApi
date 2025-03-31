package request

import (
	"fmt"
	"mime/multipart"
	"strings"
	"unicode/utf8"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	postCommand "github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UpdatePostRequest struct {
	Content  *string                `form:"content,omitempty"`
	Privacy  *consts.PrivacyLevel   `form:"privacy,omitempty" binding:"omitempty"`
	Location *string                `form:"location,omitempty"`
	MediaIDs []uint                 `form:"media_ids,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"omitempty"`
}

func ValidateUpdatePostRequest(req interface{}) error {
	dto, ok := req.(*UpdatePostRequest)
	if !ok {
		return fmt.Errorf("validate CreatePostRequest failed")
	}

	if dto.Media != nil {
		for _, fileHeader := range dto.Media {
			if err := validateMediaForUpdate(&fileHeader); err != nil {
				return err
			}
		}
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.Content, validation.By(func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid content type")
			}

			length := utf8.RuneCountInString(str)
			if length < 2 || length > 10000 {
				return fmt.Errorf("content length must be between 2 and 10000 characters, but got %d", length)
			}
			return nil
		})),
		validation.Field(&dto.Privacy, validation.In(consts.PrivacyLevels...)),
	)
}

func validateMediaForUpdate(value interface{}) error {
	if value == nil {
		return nil
	}

	fileHeader, ok := value.(*multipart.FileHeader)
	if !ok {
		return nil
	}

	if fileHeader.Size == 0 {
		return nil
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !(strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "video/")) {
		return fmt.Errorf("file must be an image or video")
	}

	// if fileHeader.Size > 10*1024*1024 {
	// 	return fmt.Errorf("file size must be less than 10M")
	// }

	return nil
}

func (req *UpdatePostRequest) ToUpdatePostCommand(
	postId *uuid.UUID,
	media []multipart.FileHeader,
) (*postCommand.UpdatePostCommand, error) {
	return &postCommand.UpdatePostCommand{
		PostId:   postId,
		Content:  req.Content,
		Privacy:  req.Privacy,
		Location: req.Location,
		MediaIDs: req.MediaIDs,
		Media:    media,
	}, nil
}
