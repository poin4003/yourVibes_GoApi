package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type LikePostCommand struct {
	UserId uuid.UUID
	PostId uuid.UUID
}

type LikePostCommandResult struct {
	Post           *common.PostResultWithLiked
	ResultCode     int
	HttpStatusCode int
}
