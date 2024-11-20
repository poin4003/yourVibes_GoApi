package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
)

type LikeCommentCommand struct {
	UserId    uuid.UUID
	CommentId uuid.UUID
}

type LikeCommentResult struct {
	Comment        *common.CommentResultWithLiked
	ResultCode     int
	HttpStatusCode int
}
