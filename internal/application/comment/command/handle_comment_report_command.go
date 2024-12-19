package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
)

type HandleCommentReportCommand struct {
	UserId uuid.UUID
}

type HandleCommentReportCommandResult struct {
	CommentReport  *common.CommentReportResult
	ResultCode     int
	HttpStatusCode int
}
