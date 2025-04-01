package request

import (
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UpdateUserRequest struct {
	FamilyName      *string              `form:"family_name,omitempty"`
	Name            *string              `form:"name,omitempty"`
	PhoneNumber     *string              `form:"phone_number,omitempty"`
	Birthday        *time.Time           `form:"birthday,omitempty"`
	Avatar          multipart.FileHeader `form:"avatar_url,omitempty" binding:"omitempty"`
	Capwall         multipart.FileHeader `form:"capwall_url,omitempty" binding:"omitempty"`
	Privacy         *consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty"`
	Biography       *string              `form:"biography,omitempty"`
	LanguageSetting *consts.Language     `form:"language_setting,omitempty" binding:"omitempty"`
}

func ValidateUpdateUserRequest(req interface{}) error {
	dto, ok := req.(*UpdateUserRequest)
	if !ok {
		return fmt.Errorf("validate UpdateUserRequest failed")
	}

	return validation.ValidateStruct(dto,
		validation.Field(&dto.FamilyName, validation.RuneLength(2, 255)),
		validation.Field(&dto.Name, validation.RuneLength(2, 255)),
		validation.Field(&dto.PhoneNumber, validation.RuneLength(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&dto.Avatar, validation.By(validateImage)),
		validation.Field(&dto.Capwall, validation.By(validateImage)),
		validation.Field(&dto.Privacy, validation.In(consts.PrivacyLevels...)),
		validation.Field(&dto.Biography, validation.RuneLength(0, 500)),
		validation.Field(&dto.LanguageSetting, validation.In(consts.Languages...)),
	)
}

func validateImage(value interface{}) error {
	if value == nil {
		return nil
	}

	fileHeader, ok := value.(multipart.FileHeader)
	if !ok {
		return nil
	}

	if fileHeader.Size == 0 {
		return nil
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("file must be an image")
	}

	// if fileHeader.Size > 10*1024*1024 {
	// 	return fmt.Errorf("file size must be less than 10MB")
	// }

	return nil
}

func (req *UpdateUserRequest) ToUpdateUserCommand(
	userId uuid.UUID,
	avatar *multipart.FileHeader,
	capwall *multipart.FileHeader,
) (*userCommand.UpdateUserCommand, error) {
	return &userCommand.UpdateUserCommand{
		UserId:          &userId,
		FamilyName:      req.FamilyName,
		Name:            req.Name,
		PhoneNumber:     req.PhoneNumber,
		Birthday:        req.Birthday,
		Avatar:          avatar,
		Capwall:         capwall,
		Privacy:         req.Privacy,
		Biography:       req.Biography,
		LanguageSetting: req.LanguageSetting,
	}, nil
}
