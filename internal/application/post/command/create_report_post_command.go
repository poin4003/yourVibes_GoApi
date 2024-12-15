package command

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type CreateReportPostCommand struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	Reason         string
}

type CreateReportPostCommandResult struct {
	PostReport     *common.PostReportResult
	ResultCode     int
	HttpStatusCode int
}
