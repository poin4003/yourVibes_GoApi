package request

import (
	"github.com/google/uuid"
	user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
	"time"
)

type UpdateUserRequest struct {
	FamilyName      *string              `form:"family_name,omitempty"`
	Name            *string              `form:"name,omitempty"`
	Email           *string              `form:"email,omitempty"`
	PhoneNumber     *string              `form:"phone_number,omitempty"`
	Birthday        *time.Time           `form:"birthday,omitempty"`
	Avatar          multipart.FileHeader `form:"avatar_url,omitempty" binding:"omitempty,file"`
	Capwall         multipart.FileHeader `form:"capwall_url,omitempty" binding:"omitempty,file"`
	Privacy         *consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Biography       *string              `form:"biography,omitempty"`
	LanguageSetting *consts.Language     `form:"language_setting,omitempty" binding:"omitempty,language_setting"`
}

func (req *UpdateUserRequest) ToUpdateUserCommand(
	userId uuid.UUID,
	avatar multipart.File,
	capwall multipart.File,
) (*user_command.UpdateUserCommand, error) {
	return &user_command.UpdateUserCommand{
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
