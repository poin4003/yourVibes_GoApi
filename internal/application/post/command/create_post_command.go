package command

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreatePostCommand struct {
	UserId   uuid.UUID
	Content  string
	Privacy  consts.PrivacyLevel
	Location string
	Media    []multipart.FileHeader
}

type CreatePostCommandResult struct {
	Post *common.PostResult
}
