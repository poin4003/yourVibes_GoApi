package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
	"time"
)

type UpdateUserCommand struct {
	UserId          *uuid.UUID
	FamilyName      *string
	Name            *string
	PhoneNumber     *string
	Birthday        *time.Time
	Avatar          multipart.File
	Capwall         multipart.File
	Privacy         *consts.PrivacyLevel
	Biography       *string
	LanguageSetting *consts.Language
}

type UpdateUserCommandResult struct {
	User           *common.UserWithSettingResult
	ResultCode     int
	HttpStatusCode int
}
