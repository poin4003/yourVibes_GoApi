package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type SharePostCommand struct {
	UserId   uuid.UUID
	PostId   uuid.UUID
	Content  string
	Privacy  consts.PrivacyLevel
	Location string
}

type SharePostCommandResult struct {
	Post           *common.PostResult
	ResultCode     int
	HttpStatusCode int
}
