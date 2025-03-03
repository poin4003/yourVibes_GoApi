package command

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type UpdatePostCommand struct {
	PostId   *uuid.UUID
	Content  *string
	Privacy  *consts.PrivacyLevel
	Location *string
	MediaIDs []uint
	Media    []multipart.FileHeader
}

type UpdatePostCommandResult struct {
	Post *common.PostResult
}
