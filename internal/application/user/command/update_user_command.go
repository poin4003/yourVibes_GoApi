package command

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UpdateUserCommand struct {
	UserId          *uuid.UUID
	FamilyName      *string
	Name            *string
	PhoneNumber     *string
	Birthday        *time.Time
	Avatar          *multipart.FileHeader
	Capwall         *multipart.FileHeader
	Privacy         *consts.PrivacyLevel
	Biography       *string
	LanguageSetting *consts.Language
}

type UpdateUserCommandResult struct {
	User *common.UserWithSettingResult
}
