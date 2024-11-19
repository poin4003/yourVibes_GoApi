package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

type CreatePostCommand struct {
	UserId   uuid.UUID
	Content  string
	Privacy  consts.PrivacyLevel
	Location string
	Media    []multipart.File
}

type CreatePostCommandResult struct {
	Post           *common.PostResult
	ResultCode     int
	HttpStatusCode int
}
